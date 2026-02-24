package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListQuestions(c *gin.Context) {
	c.JSON(200, gin.H{"questions": []interface{}{}, "total": 0})
}

func GetQuestion(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{"id": id, "question_text": "sample"})
}

func CreateQuestion(c *gin.Context) {
	c.JSON(201, gin.H{"id": uuid.New()})
}

func UpdateQuestion(c *gin.Context) {
	c.JSON(200, gin.H{"status": "updated"})
}

func DeleteQuestion(c *gin.Context) {
	c.JSON(200, gin.H{"status": "deleted"})
}

func GetProfile(c *gin.Context) {
	c.JSON(200, gin.H{
		"id":       uuid.New(),
		"username": "student1",
		"points":   150,
		"level":    5,
	})
}

func UpdateProfile(c *gin.Context) {
	c.JSON(200, gin.H{"status": "updated"})
}

func GetProgress(c *gin.Context) {
	c.JSON(200, gin.H{
		"total_questions": 100,
		"correct":         75,
		"streak":          7,
	})
}

func GetStats(c *gin.Context) {
	c.JSON(200, gin.H{
		"practice_time":    "12h 30m",
		"questions_solved": 250,
		"accuracy":         0.78,
	})
}

func GetBadges(c *gin.Context) {
	c.JSON(200, gin.H{"badges": []interface{}{
		map[string]interface{}{"code": "first_step", "name": "First Step"},
		map[string]interface{}{"code": "streak_7", "name": "Week Warrior"},
	}})
}

func StartPractice(c *gin.Context) {
	c.JSON(200, gin.H{"session_id": uuid.New()})
}

func SubmitAnswer(c *gin.Context) {
	c.JSON(200, gin.H{"correct": true})
}

func EndPractice(c *gin.Context) {
	c.JSON(200, gin.H{"score": 85})
}

func GetLeaderboard(c *gin.Context) {
	c.JSON(200, gin.H{"rankings": []interface{}{
		map[string]interface{}{"rank": 1, "username": "top_student", "points": 5000},
		map[string]interface{}{"rank": 2, "username": "second", "points": 4500},
	}})
}

func ListSubjects(c *gin.Context) {
	c.JSON(200, gin.H{"subjects": []map[string]interface{}{
		{"id": 1, "code": "math", "name_fr": "Mathématiques", "name_ar": "الرياضيات"},
		{"id": 2, "code": "pc", "name_fr": "Physique-Chimie", "name_ar": "الفيزياء والكيمياء"},
		{"id": 3, "code": "svt", "name_fr": "SVT", "name_ar": "العلوم"},
	}})
}
