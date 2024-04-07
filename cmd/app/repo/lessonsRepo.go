package repo

import (
	"errors"
	"gorm.io/gorm"
	"website-telemetry-demo/cmd/app/entities"
)

type LessonsRepo struct {
	*gorm.DB
}

func NewLessonsRepo(DB *gorm.DB) LessonsRepo {
	return LessonsRepo{DB: DB}
}

func (repo LessonsRepo) GetAllLessons() ([]entities.Lesson, error) {
	var lessons []entities.Lesson

	err := repo.DB.Omit("Description", "Likes", "Dislikes").Find(&lessons).Limit(50).Error

	return lessons, err
}

func (repo LessonsRepo) GetLessonByID(id uint64) (entities.Lesson, error) {
	var lesson entities.Lesson

	err := repo.DB.Preload("Comments").First(&lesson, id).Error
	if err != nil {
		return entities.Lesson{}, err
	} else if lesson.ID == 0 {
		return entities.Lesson{}, errors.New("lesson does not exist")
	}

	return lesson, err
}

func (repo LessonsRepo) LikeLesson(lessonID uint64) error {
	var lesson entities.Lesson

	err := repo.DB.First(&lesson, lessonID).Error
	if err != nil {
		return err
	} else if lesson.ID == 0 {
		return errors.New("lesson does not exist")
	}

	lesson.Likes++

	return repo.DB.Save(&lesson).Error
}

func (repo LessonsRepo) DislikeLesson(lessonID uint64) error {
	var lesson entities.Lesson

	err := repo.DB.First(&lesson, lessonID).Error
	if err != nil {
		return err
	} else if lesson.ID == 0 {
		return errors.New("lesson does not exist")
	}

	lesson.Dislikes++

	return repo.DB.Save(&lesson).Error
}

func (repo LessonsRepo) CommentLesson(lessonID uint64, comment entities.Comment) error {
	var lesson entities.Lesson

	err := repo.DB.First(&lesson, lessonID).Error
	if err != nil {
		return err
	} else if lesson.ID == 0 {
		return errors.New("lesson does not exist")
	}

	comment.LessonID = &lesson.ID

	return repo.DB.Create(&comment).Error
}
