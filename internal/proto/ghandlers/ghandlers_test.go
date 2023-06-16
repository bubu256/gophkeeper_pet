package ghandlers_test

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"

	// "github.com/bubu256/gophkeeper_pet/internal/proto/ghandlers"
	"github.com/bubu256/gophkeeper_pet/internal/proto/ghandlers"
	pb "github.com/bubu256/gophkeeper_pet/internal/proto/pb"
	"github.com/bubu256/gophkeeper_pet/internal/schema"
	"github.com/stretchr/testify/assert"
)

func TestConvertSchemaInfoCellToPB(t *testing.T) {
	schemaCell := &schema.InfoCell{
		ID:          1,
		DataType:    "string",
		DataSize:    10,
		Description: "Sample info cell",
		OwnerID:     1,
	}

	pbCell := ghandlers.ConvertSchemaInfoCellToPB(schemaCell)

	expectedCell := &pb.InfoCell{
		Id:          1,
		DataType:    "string",
		DataSize:    10,
		Description: "Sample info cell",
		OwnerId:     1,
	}

	assert.Equal(t, expectedCell, pbCell)
}

func TestConvertSchemaMemoryCellToPB(t *testing.T) {
	schemaCell := &schema.MemoryCell{
		InfoCell: &schema.InfoCell{
			ID:          1,
			DataType:    "string",
			DataSize:    10,
			Description: "Sample info cell",
			OwnerID:     1,
		},
		ID:            1,
		Encrypted:     true,
		KeyValuePairs: map[string]string{"key": "value"},
		BinaryData:    []byte{0x01, 0x02, 0x03},
		FileName:      "data.txt",
	}

	pbCell := ghandlers.ConvertSchemaMemoryCellToPB(schemaCell)

	expectedCell := &pb.MemoryCell{
		Id: 1,
		Info: &pb.InfoCell{
			Id:          1,
			DataType:    "string",
			DataSize:    10,
			Description: "Sample info cell",
			OwnerId:     1,
		},
		Encrypted:     true,
		KeyValuePairs: map[string]string{"key": "value"},
		BinaryData:    []byte{0x01, 0x02, 0x03},
		FileName:      "data.txt",
	}

	assert.Equal(t, expectedCell, pbCell)
}

func TestGetTokenFromContext(t *testing.T) {
	md := map[string]string{
		"token": "example_token",
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(md))

	token, ok := ghandlers.GetTokenFromContext(ctx)

	assert.True(t, ok)
	assert.Equal(t, "example_token", token)
}

func TestConvertPBInfoCellToSchema(t *testing.T) {
	pbCell := &pb.InfoCell{
		Id:          1,
		DataType:    "string",
		DataSize:    10,
		Description: "Sample info cell",
		OwnerId:     1,
	}

	schemaCell := ghandlers.ConvertPBInfoCellToSchema(pbCell)

	expectedCell := &schema.InfoCell{
		ID:          1,
		DataType:    "string",
		DataSize:    10,
		Description: "Sample info cell",
		OwnerID:     1,
	}

	assert.Equal(t, expectedCell, schemaCell)
}

func TestConvertPBMemoryCellToSchema(t *testing.T) {
	pbCell := &pb.MemoryCell{
		Id: 1,
		Info: &pb.InfoCell{
			Id:          1,
			DataType:    "string",
			DataSize:    10,
			Description: "Sample info cell",
			OwnerId:     1,
		},
		Encrypted:     true,
		KeyValuePairs: map[string]string{"key": "value"},
		BinaryData:    []byte{0x01, 0x02, 0x03},
		FileName:      "data.txt",
	}

	schemaCell := ghandlers.ConvertPBMemoryCellToSchema(pbCell)

	expectedCell := &schema.MemoryCell{
		InfoCell: &schema.InfoCell{
			ID:          1,
			DataType:    "string",
			DataSize:    10,
			Description: "Sample info cell",
			OwnerID:     1,
		},
		ID:            1,
		InfoID:        1,
		Encrypted:     true,
		KeyValuePairs: map[string]string{"key": "value"},
		BinaryData:    []byte{0x01, 0x02, 0x03},
		FileName:      "data.txt",
	}

	assert.Equal(t, expectedCell, schemaCell)
}
