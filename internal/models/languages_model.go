package models

const (
	LanguageEnglish     = "en" // English
	LanguageVietnamese  = "vi" // Tiếng Việt
	// LanguageFrench      = "fr" // Français
	// LanguageSpanish     = "es" // Español
	LanguageJapanese    = "ja" // 日本語 (Nihongo)
	LanguageKorean      = "ko" // 한국어 (Hangul)
	LanguageChineseCN   = "zh-CN" // 中文 (简体, 中国)
	// LanguageChineseTW   = "zh-TW" // 中文 (繁體, 台灣)
	// LanguageGerman      = "de" // Deutsch
	AutoLanguageDetection = "auto"
)

var ValidLanguages = []string{
	LanguageEnglish,
	LanguageVietnamese,
	// LanguageFrench,
	// LanguageSpanish,
	LanguageJapanese,
	LanguageKorean,
	LanguageChineseCN,
	// LanguageChineseTW,
	// LanguageGerman,
	AutoLanguageDetection,
}