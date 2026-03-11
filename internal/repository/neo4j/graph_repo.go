package neo4j

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphRepo struct {
	driver neo4j.Driver
}

func NewGraphRepo(driver neo4j.Driver) *GraphRepo {
	return &GraphRepo{driver: driver}
}

func (r *GraphRepo) CreateStudentNode(studentID string, name string) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MERGE (s:Student {id: $studentID})
		SET s.name = $name
	`

	_, err := session.Run(query, map[string]interface{}{
		"studentID": studentID,
		"name":      name,
	})

	return err
}

func (r *GraphRepo) CreateDimensionNode(name string) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MERGE (d:Dimension {name: $name})
	`

	_, err := session.Run(query, map[string]interface{}{
		"name": name,
	})

	return err
}

func (r *GraphRepo) UpdateKnowledgeGraph(studentID string, dimension string, score float64) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (s:Student {id: $studentID})
		MATCH (d:Dimension {name: $dimension})
		MERGE (s)-[r:HAS_SCORE]->(d)
		SET r.score = $score, r.updatedAt = timestamp()
	`

	_, err := session.Run(query, map[string]interface{}{
		"studentID": studentID,
		"dimension": dimension,
		"score":     score,
	})

	return err
}

func (r *GraphRepo) GetStudentScores(studentID string) (map[string]float64, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (s:Student {id: $studentID})-[r:HAS_SCORE]->(d:Dimension)
		RETURN d.name AS dimension, r.score AS score
	`

	result, err := session.Run(query, map[string]interface{}{
		"studentID": studentID,
	})
	if err != nil {
		return nil, err
	}

	scores := make(map[string]float64)
	for result.Next() {
		record := result.Record()
		dimension, _ := record.Get("dimension")
		score, _ := record.Get("score")
		if dimension != nil && score != nil {
			scores[dimension.(string)] = score.(float64)
		}
	}

	return scores, nil
}
