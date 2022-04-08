package main

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	//"github.com/qinains/fastergoding"
)

// Blockstructure
type BlockData struct {
	Title             string
	Description       string
	Owners            []string
	Problem           string
	Domain            []string
	Technologies_used []string
	Viewing_price     float32
	Ownership_price   float32
	Pricing_history   []float32
}

// Blockchain structure
type Block struct {
	Data        BlockData
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}

// For proposedIdeas cache
type ProposedIdea struct {
	BlockData BlockData
	SimIdea   string
	SimScore  float32
}

type User struct {
	Username    string
	Password    string
	Balance     float32
	Email       string
	PhoneNumber string
}

// expect head pointer which will now point to the tail of the blockchain
func InsertBlock(chainHead *Block, newBlock BlockData) *Block {
	if chainHead != nil {
		newBlock := &Block{
			Data:        newBlock,
			PrevHash:    "ph2",
			CurrentHash: "ch2",
			PrevPointer: chainHead,
		}
		chainHead = newBlock
	} else {
		chainHead = &Block{
			Data:        newBlock,
			PrevHash:    "ph1",
			CurrentHash: "ch1",
			PrevPointer: nil,
		}
	}
	return chainHead
}

func PrintChain(chainHead *Block) {
	cursor := chainHead
	for {
		if cursor != nil {
			println(cursor.Data.Title)
			println(cursor.Data.Viewing_price)
			println(cursor.Data.Description)
			cursor = cursor.PrevPointer
		} else {
			return
		}
	}
}

// GEt full chain
func GetChain(chain *Block) []BlockData {
	blockchain := []BlockData{}
	cursor := chain
	for {
		if cursor != nil {
			blockchain = append(blockchain, cursor.Data)
			cursor = cursor.PrevPointer
		} else {
			return blockchain
		}
	}
}

func splice(proposedIdea []ProposedIdea, index int) []ProposedIdea {
	return append(proposedIdea[:index], proposedIdea[index+1:]...)
}

func findIndexOfProposal(proposalToFind ProposedIdea, proposedIdeas []ProposedIdea) int {
	for i := 0; i < len(proposedIdeas); i++ {
		if proposalToFind.BlockData.Title == proposedIdeas[i].BlockData.Title &&
			proposalToFind.BlockData.Description == proposedIdeas[i].BlockData.Description &&
			proposalToFind.BlockData.Viewing_price == proposedIdeas[i].BlockData.Viewing_price &&
			proposalToFind.BlockData.Ownership_price == proposedIdeas[i].BlockData.Ownership_price {
			return i
		}
	}
	return -1
}

func main() {

	var ChainHead *Block

	Proposals := []ProposedIdea{}

	// test block
	blockd := BlockData{
		Title:             "ME",
		Description:       "IM",
		Problem:           "Ima mess",
		Owners:            []string{},
		Domain:            []string{},
		Technologies_used: []string{},
		Viewing_price:     1.4,
		Ownership_price:   0,
		Pricing_history:   []float32{},
	}
	// test adding data
	ChainHead = InsertBlock(ChainHead, blockd)
	ChainHead = InsertBlock(ChainHead, blockd)
	ChainHead = InsertBlock(ChainHead, blockd)
	// PrintChain(ChainHead)

	// -------------SERVER SIDE PROGRAM----------
	//fastergoding.Run()
	server := fiber.New()
	server.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to GOBE")
	})

	// Add Idea to Chain (Approved idea)
	server.Post("/addidea", func(c *fiber.Ctx) error {
		var blockData BlockData
		if json.Unmarshal(c.Body(), &blockData) == nil {
			ChainHead = InsertBlock(ChainHead, blockData)
			return c.SendStatus(200)
		}
		return c.SendStatus(500)
	})

	// Get the block chain
	server.Get("/getchain", func(c *fiber.Ctx) error {
		return c.JSON(GetChain(ChainHead))
	})

	// Recieve the proposals from PY. PY will calculate the data and send on this endpoint
	server.Post("/proposeidea", func(c *fiber.Ctx) error {

		print("muhhahha")
		var propsedIdea ProposedIdea
		if json.Unmarshal(c.Body(), &propsedIdea) == nil {
			return c.SendStatus(500)
		}
		println(propsedIdea.SimScore)
		fmt.Printf("%f\n", propsedIdea.SimScore)
		Proposals = append(Proposals, propsedIdea)
		return c.SendStatus(200)
	})

	// FE will see this
	server.Get("/getproposals", func(c *fiber.Ctx) error {
		return c.JSON(Proposals)
	})

	// approve an idea from proposals
	server.Post("/approveproposal", func(c *fiber.Ctx) error {
		var recievedProposedIdea ProposedIdea
		json.Unmarshal(c.Body(), &recievedProposedIdea)
		extractedDataFromProposal := recievedProposedIdea.BlockData
		Proposals = splice(Proposals, findIndexOfProposal(recievedProposedIdea, Proposals))
		ChainHead = InsertBlock(ChainHead, extractedDataFromProposal)
		return c.JSON(GetChain(ChainHead))
	})
	server.Post("/disapproveproposal", func(c *fiber.Ctx) error {
		var recievedProposedIdea ProposedIdea
		json.Unmarshal(c.Body(), &recievedProposedIdea)
		Proposals = splice(Proposals, findIndexOfProposal(recievedProposedIdea, Proposals))
		return c.JSON(GetChain(ChainHead))
	})

	// TODO
	server.Post("/registeruser", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	server.Post("/updateuser", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	server.Post("/loginuser", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	server.Post("/deleteuser", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	server.Post("/addcoinstowallet", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	server.Post("/addtoauction", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	server.Post("/bidonidea", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	server.Post("/sellideato", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	server.Listen(":8081")

}
