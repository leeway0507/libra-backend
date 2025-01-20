package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"libra-backend/db/sqlc"
	"libra-backend/pkg/pb"
	"libra-backend/pkg/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/proto"
)

const maxTokenSize = 7000

type QueryReq struct {
	query     *sqlc.Queries
	openAIKey string
	dataPath  string
	batchId   string
}

type ReqBody struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

type RespBody struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int32 `json:"prompt_tokens"`
		TotalTokens  int32 `json:"total_tokens"`
	} `json:"usage"`
}

type QueryResp struct {
	Query     string
	Embedding []float32
}

func New(query *sqlc.Queries, openAIKey string) *QueryReq {
	path := CreateKeywordFolder()
	log.Printf("path: %#+v\n", path)
	return &QueryReq{
		query,
		openAIKey,
		path,
		"batchId",
	}
}
func CreateKeywordFolder() string {
	path := filepath.Join(utils.GetCurrentFolderDir(), "query")
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Printf("err: %#+v\n", err)
	}
	return utils.GetCurrentFolderDir()
}

func (Q *QueryReq) DBQuery() *sqlc.Queries {
	return Q.query
}

func (Q *QueryReq) RequestQueryEmbedding(query string) (*QueryResp, error) {
	queryEmbedding, err := Q.LoadEmbeddingQuery(query)
	if err == nil {
		log.Printf("Query Found: %#+v\n", query)
		return &QueryResp{
			Query:     queryEmbedding.Query,
			Embedding: queryEmbedding.Embedding,
		}, nil
	}
	runes := []rune(query)

	reqBody := &ReqBody{
		Input: string(runes[0:min(len(runes), maxTokenSize)]),
		Model: "text-embedding-3-small",
	}

	reqBodyByte, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := "https://api.openai.com/v1/embeddings"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyByte))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Q.openAIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var openAIresp RespBody
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	json.Unmarshal(body, &openAIresp)
	path := filepath.Join(Q.dataPath, "query", "temp"+".txt")
	temp, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	temp.Write(body)

	return &QueryResp{
		Query:     query,
		Embedding: openAIresp.Data[0].Embedding,
	}, nil
}
func (Q *QueryReq) SaveQueryEmbedding(queryEmbedding *pb.QueryEmbedding) {
	path := filepath.Join(Q.dataPath, "query", queryEmbedding.Query+".pb")

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		panic("fail to make directories")
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := proto.Marshal(queryEmbedding)
	if err != nil {
		panic(err)
	}
	file.Write(b)
}
func (Q *QueryReq) LoadEmbeddingQuery(query string) (*pb.QueryEmbedding, error) {
	path := filepath.Join(Q.dataPath, "query", query+".pb")
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	b := Q.LoadFile(path)
	embeddingVector := &pb.QueryEmbedding{}
	err := proto.Unmarshal(b, embeddingVector)
	if err != nil {
		return nil, err
	}

	return embeddingVector, nil
}

func (R *QueryReq) LoadFile(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file: %v", err))
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("Failed to read all: %v", err))
	}
	return b
}
