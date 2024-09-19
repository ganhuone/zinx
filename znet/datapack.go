package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/ganhuone/zinx/utils"
	"github.com/ganhuone/zinx/ziface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	// unit32 + unit32 = len + id = 8
	return 8
}

func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen())
	if err != nil {
		return nil, err
	}

	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (d *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)

	msg := &Message{}

	err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	if msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data reve")
	}

	err = binary.Read(dataBuff, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}

	return msg, nil

}
