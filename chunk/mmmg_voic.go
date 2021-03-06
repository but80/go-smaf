package chunk

import (
	"io"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/but80/go-smaf.v1/enums"
	"gopkg.in/but80/go-smaf.v1/subtypes"
	"gopkg.in/but80/go-smaf.v1/util"
)

type MMMGVoiceChunk struct {
	*ChunkHeader
	SubChunks []Chunk `json:"sub_chunks"`
}

func (c *MMMGVoiceChunk) GetExclusives() []*subtypes.Exclusive {
	result := []*subtypes.Exclusive{}
	c.Traverse(func(s Chunk) {
		switch exvo := s.(type) {
		case *MMMGEXVOChunk:
			result = append(result, exvo.Exclusive)
		}
	})
	return result
}

func (c *MMMGVoiceChunk) Traverse(fn func(Chunk)) {
	fn(c)
	for _, sub := range c.SubChunks {
		sub.Traverse(fn)
	}
}

func (c *MMMGVoiceChunk) String() string {
	result := "MMMGVoiceChunk: " + c.ChunkHeader.String()
	sub := []string{}
	for _, chunk := range c.SubChunks {
		sub = append(sub, chunk.String())
	}
	return result + "\n" + util.Indent(strings.Join(sub, "\n"), "\t")
}

// Read は、バイト列を読み取ってパースした結果をこの構造体に格納します。
func (c *MMMGVoiceChunk) Read(rdr io.Reader) error {
	rest := int(c.ChunkHeader.Size)
	for 8 <= rest {
		var hdr ChunkHeader
		err := hdr.Read(rdr, &rest)
		if err != nil {
			return errors.WithStack(err)
		}
		sub, err := hdr.CreateChunk(rdr, enums.ScoreTrackFormatType_Default)
		if err != nil {
			return errors.WithStack(err)
		}
		c.SubChunks = append(c.SubChunks, sub)
	}
	return nil
}
