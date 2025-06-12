package classifier

type Classifier interface {
	Classify(text string) (string, error)
}
