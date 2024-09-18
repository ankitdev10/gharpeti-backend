// handlers/user.go
package handlers

import (
  "errors"
  "fmt"
  "gharpeti/cmd/db"
  "gharpeti/dto"
  "gharpeti/models"
  "gharpeti/utils"
  "github.com/labstack/echo/v4"
  "golang.org/x/crypto/bcrypt"
  "gorm.io/gorm"
  "net/http"
)

func CreateUser(c echo.Context) error {
  dto := c.Get("dto").(*dto.CreateUserDTO)

  result := db.DB.Where("email = ?", dto.Email).First(&models.User{})

  if result.RowsAffected > 0 {
    return utils.SendError(c, 400, "User already exists")
  }

  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
  if err != nil {
    fmt.Println(err)
    return utils.SendError(c, http.StatusInternalServerError, "Error hashing password")
  }

  newUser := models.User{
    FullName: dto.FullName,
    Email:    dto.Email,
    Password: string(hashedPassword),
    Location: dto.Location,
    Phone:    dto.Phone,
    Type:     dto.Type,
  }

  newUser.Password = string(hashedPassword)

  if err := db.DB.Create(&newUser).Error; err != nil {
    fmt.Println(err)
    return utils.SendError(c, http.StatusInternalServerError, "Error creating User")
  }

  newUser.Password = ""

  return utils.SendSuccessResponse(c, http.StatusCreated, "User created", newUser)
}
func GetUser(c echo.Context) error {
  var users []models.User
  result := db.DB.Find(&users)
  if result.Error != nil {
    fmt.Println(result.Error)
    return utils.SendError(c, http.StatusInternalServerError, "Error fetching users")
  }
  return utils.SendSuccessResponse(c, http.StatusOK, "Users fetched", users)
}

func GetOneUser(c echo.Context) error {
  var user models.User
  id := c.Param("id")
  db := db.DB
  result := db.First(&user, id)

  if result.Error != nil {

    if !errors.Is(result.Error, gorm.ErrRecordNotFound) {

      return utils.SendError(c, http.StatusInternalServerError, "Error fetching user")
    } else {
      return utils.SendError(c, http.StatusNotFound, "User not found")
    }
  }
  user.Password = ""
  return utils.SendSuccessResponse(c, http.StatusOK, "User fetched", user)
}

func UpdateUser(c echo.Context) error {
  id := c.Param("id")
  dto := c.Get("dto").(*dto.UpdateUserDTO)
  db := db.DB
  var existingUser models.User

  findUser := db.Find(&existingUser, id)
  if findUser.Error != nil {
    if errors.Is(findUser.Error, gorm.ErrRecordNotFound) {
      return utils.SendError(c, http.StatusNotFound, "User not found")
    }
  }

  updatedUser := existingUser

  if err := c.Bind(&dto); err != nil {
    return utils.SendError(c, http.StatusInternalServerError, "Can not bind request body")
  }

  updateResult := db.Model(&existingUser).Updates(&updatedUser)
  if updateResult.Error != nil {
    return c.JSON(http.StatusInternalServerError, map[string]interface{}{
      "error": "Internal server error",
    })
  }
  return utils.SendSuccessResponse(c, http.StatusOK, "User updated", existingUser)
}

func Me(c echo.Context) error {
  user, ok := c.Get("user").(models.User)
  if !ok {
    return utils.SendError(c, http.StatusUnauthorized, "Unauthorized")
  }

  user.Password = ""
  user.Token = ""
  return utils.SendSuccessResponse(c, http.StatusOK, "User fetched", user)
}

func Logout(c echo.Context) error {
  user, ok := c.Get("user").(models.User)

  if !ok {
    return utils.SendError(c, http.StatusUnauthorized, "Unauthorized")
  }

  db.DB.Model(&user).Update("token", "")

  return utils.SendSuccessResponse(c, http.StatusOK, "User logged out", "Sucess")
}
