package services

import (
	"models"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"log"
	"golang.org/x/crypto/bcrypt"
	"utils"
	"repository"
	"controllers/viewmodels"
	"strconv"
	"github.com/dustin/go-humanize"
)

type UserService struct {
	repo *repository.Repository
	*RoleService
}

func NewUserService() *UserService {
	return &UserService{
		repo: repository.NewRepo(),
		RoleService: NewRoleService(),
	}
}

func (s *UserService) GetSecret(key []byte) ([]byte, error) {
	return nil, nil
}

func (s *UserService) UserExist(email string) bool {
	user, err := s.GetUserByEmail(email)
	return user != nil || err == nil
}

func (s *UserService) CreateUser(user *models.NewUser) (*models.User, error) {
	if s.UserExist(user.Email) {
		return nil, fmt.Errorf("User already exist: %s", user.Email)
	}

	u := models.User{
		FirstName: bluemonday.StrictPolicy().Sanitize(user.FirstName),
		LastName:  bluemonday.StrictPolicy().Sanitize(user.LastName),
		Email:     user.Email,
		Password: user.Password,
	}

	if user.CreatedBy != nil {
		u.CreatedByID = utils.ToUInt32(user.CreatedBy.ID)
	}

	role, err := s.GetRole(user.Role)
	if err != nil {
		return nil, err
	}

	hashed, err := encryptPassword(user.Password)
	if err != nil {
		return nil, err
	}

	u.RoleID = utils.ToUInt32(role.ID)
	u.Password = hashed
	tx := s.repo.DB().Begin()
	if err := tx.Save(&u).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating User: %v", err)
		return nil, err
	}

	userRole := repository.UserRole{
		UserID:utils.ToUInt32(u.ID),
		RoleID:utils.ToUInt32(role.ID)}
	if err := tx.Save(&userRole).Error; err != nil {
		tx.Rollback()
		log.Printf("Error saving UserRole: %v", err)
		return nil, err
	}

	tx.Commit()

	return &u, nil
}

func (s *UserService) GetUserByID(id uint32) (*models.User, error) {
	var user models.User
	if err := s.repo.FindByID(id, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUsers(query models.Query) ([]*models.User, error) {
	var users []*models.User
	if err := s.repo.FindByQuery(query, &users); err != nil {
		log.Printf("Error finding Users with Query %+v: %v", query, err)
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetUsersByQuery(query models.Query) ([]*viewmodels.User, error) {
	sql := `SELECT
			  u.id                                          AS user_id
              ,CONCAT_WS(' ', u.first_name, u.last_name)  AS full_name
			  ,u.nick_name
			  ,u.email
			  ,u.website
			  ,u.bio
			  ,p.id                                         AS profile_pic_id
			  ,p.url                                        AS profile_pic_url
			  ,r.id                                         AS role_id
			  ,r.name                                       AS role_name
			  ,u2.id                                        AS created_by_id
			  ,CONCAT_WS(' ', u2.first_name, u2.last_name)  AS created_by_name
			  ,(SELECT count(*)
			  FROM pages AS p
			  WHERE p.author_id = u.id)                     AS num_pages
			  ,(SELECT count(*)
			  FROM articles AS a
			  WHERE a.author_id = u.id)                     AS num_articles
			  ,(SELECT count(*)
			  FROM categories AS c
			  WHERE c.author_id = u.id)                     AS num_categories
			  ,(SELECT count(*)
			  FROM taxonomies AS t
			  WHERE t.author_id = u.id)                     AS num_taxonomies
			  ,u.created_at
			  ,u.updated_at
			FROM users AS u
			  INNER JOIN roles AS r
			    ON u.role_id = r.id
			  LEFT JOIN pictures AS p
			  ON u.profile_picture_id IS NOT NULL
			     AND u.profile_picture_id = p.id
			  LEFT JOIN users AS u2
			  ON u.created_by_id = u2.id
			WHERE u.deleted_at is NULL`

	var result []*repository.UserResult
	err := s.repo.DB().Raw(sql).
		Limit(query.Total).
		Offset(query.Offset).
		Order(fmt.Sprintf("u.%s", query.Sort)).
		Scan(&result).Error
	if err != nil {
		log.Printf("Error finding Users with Query %+v: %v", query, err)
		return nil, err
	}
	users := []*viewmodels.User{}
	for _, u := range result {
		user := &viewmodels.User{}
		user.ID = strconv.Itoa(int(u.UserID))
		user.FullName = u.FullName
		user.NickName = u.NickName
		user.Email = u.Email
		user.Website = u.Website
		user.Biography = u.Biography
		user.ProfilePicture.ID = strconv.Itoa(int(u.ProfilePicID))
		user.ProfilePicture.Url = u.ProfilePicUrl
		user.Role.ID = strconv.Itoa(int(u.RoleID))
		user.Role.Name = u.RoleName
		user.CreatedBy.ID = strconv.Itoa(int(u.CreatedByID))
		user.CreatedBy.FullName = u.CreatedByName
		user.NumPages = u.NumPages
		user.NumArticles = u.NumArticles
		user.NumCategories = u.NumCategories
		user.NumTaxonomies = u.NumTaxonomies
		user.CreatedAt = humanize.Time(u.CreatedAt)
		user.UpdatedAt = humanize.Time(u.UpdatedAt)

		users = append(users, user)
	}

	return users, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.repo.FindOne(&models.User{Email: email}, &user); err != nil {
		log.Printf("Error finding User for email %s : %v", email, err)
		return nil, err
	}

	return &user, nil
}
/*
func (s *UserService) UpdateUser(user *models.UpdateUser) error {
	user.UpdatedAt = time.Now()
	user.LastName = bluemonday.StrictPolicy().Sanitize(user.LastName)
	user.Bio = bluemonday.UGCPolicy().Sanitize(user.Bio)
	if len(user.Password) > 0 {
		hashed, err := encryptPassword(user.Password)
		if err != nil {
			return err
		}

		user.Password = hashed
	}

	if err := s.repo.Update(user.Id, user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) DeleteUser(id string) error {
	user, err := s.GetUser(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(user.ID); err != nil {
		return err
	}

	*/
/*if err := NewArticleService().DeleteArticles(user.Articles); err != nil {
		//s.UpdateUser(user)
		return err
	}*//*


	return nil
}
*/

func (s *UserService) VerifyLogin(login models.Login) (*models.User, error) {
	user, err := s.GetUserByEmail(login.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, login.Password); err != nil {
		log.Printf("Error password %s doesn't match: %v", login.Password, err)
		return nil, err
	}

	return user, nil
}

func encryptPassword(p []byte) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating bcrypt hash: %v", err)
		return nil, err
	}

	return hashed, nil
}
