package data

import "mngr/models"

type AiDto interface {
	GetAiType() int
	GetAiClip() *AiClip
	GetImageFileName() string
	GetPredClsIdx() int
	GetPredClsName() string
	GetPredScore() float32
	GetId() string
}

func (o *OdDto) GetAiType() int {
	return models.Od
}

func (o *OdDto) GetAiClip() *AiClip {
	return o.AiClip
}

func (o *OdDto) GetImageFileName() string {
	return o.ImageFileName
}

func (o *OdDto) GetPredClsIdx() int {
	return o.DetectedObject.PredClsIdx
}

func (o *OdDto) GetPredClsName() string {
	return o.DetectedObject.PredClsName
}

func (o *OdDto) GetPredScore() float32 {
	return o.DetectedObject.PredScore
}

func (o *OdDto) GetId() string {
	return o.Id
}

func (f *FrDto) GetAiType() int {
	return models.Fr
}

func (f *FrDto) GetAiClip() *AiClip {
	return f.AiClip
}

func (f *FrDto) GetImageFileName() string {
	return f.ImageFileName
}

func (f *FrDto) GetPredClsIdx() int {
	return f.DetectedFace.PredClsIdx
}

func (f *FrDto) GetPredClsName() string {
	return f.DetectedFace.PredClsName
}

func (f *FrDto) GetPredScore() float32 {
	return f.DetectedFace.PredScore
}

func (f *FrDto) GetId() string {
	return f.Id
}

func (a *AlprDto) GetAiType() int {
	return models.Alpr
}

func (a *AlprDto) GetAiClip() *AiClip {
	return a.AiClip
}

func (a *AlprDto) GetImageFileName() string {
	return a.ImageFileName
}

func (a *AlprDto) GetPredClsIdx() int {
	return 0
}

func (a *AlprDto) GetPredClsName() string {
	return a.DetectedPlate.Plate
}

func (a *AlprDto) GetPredScore() float32 {
	return float32(a.DetectedPlate.Confidence)
}

func (a *AlprDto) GetId() string {
	return a.Id
}
