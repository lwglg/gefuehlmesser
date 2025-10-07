package sentiment

const (
	HashtagRegexString          string = "^#[a-zA-Z-а-яА-ЯÀ-ÖØ-öø-ʸ0-9(_)]{1,}$"
	MessageIdRegexString        string = "^msg_[a-z0-9_]{3,}$"
	UserIdRegexString           string = "^user_[a-z0-9_]{3,}$"
	TokenizationRegexString     string = "(?:#\\w+(?:-\\w+)*)|\b\\w+\b"
	TimeStampRFC3339RegexString string = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
	PunctuationRegexString      string = "[\\.,!\\?;:\"\\(\\)\\[\\]{}…]"
)

var (
	PositiveWords = []string{"bom", "ótimo", "adorei", "excelente", "maravilhoso", "perfeito", "gostei"}
	NegativeWords = []string{"ruim", "péssimo", "odiei", "terrível", "horrível", "decepcionante"}
	Intensifiers  = []string{"muito", "super", "extremamente"}
	Negations     = []string{"não", "nunca", "jamais"}
)
