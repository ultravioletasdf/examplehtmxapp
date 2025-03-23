package routes

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"examplehtmxapp/frontend"
	sqlc "examplehtmxapp/sql"
	"examplehtmxapp/utils"
	"fmt"
	"strings"
	"time"

	_snowflake "github.com/bwmarrin/snowflake"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func onboarding(c *fiber.Ctx) error {
	session := c.Cookies("session")
	user, err := executor.GetUserFromSession(ctx, session)
	if err == sql.ErrNoRows {
		return c.Redirect("/sign/in?toast=Something went wrong... Please try signing in again", fiber.StatusTemporaryRedirect)
	}
	if user.Verified == 1 {
		return c.Redirect("/?toast=You are already verified", fiber.StatusTemporaryRedirect)
	}
	if err != nil {
		return c.SendString(err.Error())
	}
	fmt.Println(user)
	return Render(c, frontend.Onboarding())
}
func putOnboarding(c *fiber.Ctx) error {
	session := c.Cookies("session")

	pin1 := c.FormValue("pin-1")
	pin2 := c.FormValue("pin-2")
	pin3 := c.FormValue("pin-3")
	pin4 := c.FormValue("pin-4")
	pin5 := c.FormValue("pin-5")
	pin6 := c.FormValue("pin-6")
	pin := fmt.Sprintf("%s%s%s%s%s%s", pin1, pin2, pin3, pin4, pin5, pin6)

	user, err := executor.GetUserFromSession(ctx, session)
	if err == sql.ErrNoRows {
		c.Set("HX-Redirect", "/sign/in?toast=Something went wrong... Please try signing in again")
		return c.SendStatus(200)
	}
	if err != nil {
		fmt.Println(err)
		c.Set("HX-Redirect", "/onboarding?toast=Something went wrong... Please try again")
		return c.SendStatus(200)
	}

	if !user.VerifyCode.Valid || user.VerifyCode.String != pin {
		c.Set("HX-Redirect", "/onboarding?toast=Incorrect code, please try again")
		return c.SendStatus(200)
	}
	err = executor.VerifyUser(ctx, user.ID)
	if err != nil {
		fmt.Println(err)
		c.Set("HX-Redirect", "/onboarding?toast=Something went wrong... Please try again")
		return c.SendStatus(200)
	}
	c.Set("HX-Redirect", "/")
	return c.SendStatus(200)
}
func signIn(c *fiber.Ctx) error {
	if session := c.Cookies("session"); session != "" {
		c.Cookie(&fiber.Cookie{Name: "session", Path: "/", SameSite: "strict", HTTPOnly: true, Value: "", MaxAge: -1})
		executor.DeleteSession(ctx, c.Cookies("session"))
	}
	return Render(c, frontend.SignIn())
}

func postSignIn(c *fiber.Ctx) error {
	email := strings.TrimSpace(c.FormValue("email"))
	password := strings.TrimSpace(c.FormValue("password"))
	if email == "" || password == "" {
		return Render(c, frontend.SoftError("Error: Invalid Request Body"))
	}
	if len(password) < 8 || len(password) > 72 {
		return Render(c, frontend.SoftError("Your password must be between 8 and 72 characters"))
	}
	user, err := executor.GetUserByEmail(ctx, email)
	if err != nil {
		return Render(c, frontend.SoftError("There is no account with that email"))
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return Render(c, frontend.SoftError("Incorrect Password"))
	}
	if err != nil {
		fmt.Printf("Bcrypt Error: %v\n", err)
		return Render(c, frontend.SoftError("Something went wrong... Please try signing in again"))
	}

	err = createSession(c, _snowflake.ID(user.ID))
	if err != nil {
		return err
	}
	if user.Verified != 1 {
		c.Set("HX-Redirect", "/onboarding")
		var code sql.NullString
		code.String = utils.RandomCode()
		code.Valid = true
		err = executor.SetVerification(ctx, sqlc.SetVerificationParams{VerifyCode: code, ID: user.ID})
		err := utils.SendVerification(cfg, email, code.String)
		if err != nil {
			fmt.Println(err.Error())
			return Render(c, frontend.SoftError("Internal Error: "+err.Error()))
		}
		err = utils.SendVerification(cfg, email, code.String)
		if err != nil {
			fmt.Println(err.Error())
			return Render(c, frontend.SoftError("Internal Error: "+err.Error()))
		}
		return c.SendStatus(200)
	}
	c.Set("HX-Redirect", "/")
	return c.SendStatus(200)
}
func signUp(c *fiber.Ctx) error {
	return Render(c, frontend.SignUp())
}
func postSignUp(c *fiber.Ctx) error {
	email := strings.TrimSpace(c.FormValue("email"))
	password := strings.TrimSpace(c.FormValue("password"))
	if email == "" || password == "" {
		return Render(c, frontend.SoftError("Error: Invalid Request Body"))
	}
	if len(password) < 8 || len(password) > 72 {
		return Render(c, frontend.SoftError("Your password must be between 8 and 72 characters"))
	}
	user, err := executor.GetUserByEmail(ctx, email)
	if err != sql.ErrNoRows && err != nil {
		fmt.Println(err.Error())
		return Render(c, frontend.SoftError("Internal Error: "+err.Error()))
	}
	fmt.Println(user)
	if user.Email != "" {
		return Render(c, frontend.SoftError("That email is already being used"))
	}

	id := snowflake.Generate()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
		return Render(c, frontend.SoftError("Internal Error: "+err.Error()))
	}
	err = executor.CreateUser(ctx, sqlc.CreateUserParams{ID: id.Int64(), Email: email, Password: passwordHash})
	if err != nil {
		fmt.Println(err.Error())
		return Render(c, frontend.SoftError("Internal Error: "+err.Error()))
	}

	err = createSession(c, id)
	if err != nil {
		return err
	}

	var code sql.NullString
	code.String = utils.RandomCode()
	code.Valid = true

	err = executor.SetVerification(ctx, sqlc.SetVerificationParams{VerifyCode: code, ID: id.Int64()})
	if err != nil {
		fmt.Println(err.Error())
		return Render(c, frontend.SoftError("Internal Error: "+err.Error()))
	}
	fmt.Println("Set verification", code, user.ID)

	err = utils.SendVerification(cfg, email, code.String)
	if err != nil {
		fmt.Println(err.Error())
		return Render(c, frontend.SoftError("Internal Error: "+err.Error()))
	}

	c.Set("HX-Redirect", "/onboarding")
	return c.SendStatus(200)
}
func signOut(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{Name: "session", Path: "/", SameSite: "strict", HTTPOnly: true, Value: "", MaxAge: -1})
	executor.DeleteSession(ctx, c.Cookies("session"))
	return c.Redirect("/", fiber.StatusTemporaryRedirect)
}

// c.Cookie(&fiber.Cookie{Name: "session", Path: "/", SameSite: "strict", HTTPOnly: true, Value: "", MaxAge: -1})
func createToken(size int) string {
	bytes := make([]byte, size)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}

func createSession(c *fiber.Ctx, id _snowflake.ID) error {
	token := createToken(16)
	expireAt := time.Now().Add(time.Hour * 24 * 365)
	err := executor.CreateSession(ctx, sqlc.CreateSessionParams{UserID: id.Int64(), Token: token, ExpireAt: expireAt.Unix()})
	if err != nil {
		fmt.Println(id, err.Error())
		return Render(c, frontend.SoftError("Internal Error: "+err.Error()))
	}
	c.Cookie(&fiber.Cookie{Name: "session", Value: token, Path: "/", Expires: expireAt, SameSite: "strict", HTTPOnly: true})
	return nil
}
