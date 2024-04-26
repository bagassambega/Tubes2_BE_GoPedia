package main

import "strings"

// Jika mengandung salah satu dari identifier seperti File: pada URL, return true
func checkIgnoredLink(url string) bool {
	ignoredLinks := [...]string{"/File:", "/Special:", "/Template:", "/Template_page:", "/Help:", "/Category:", "Special:", "/Wikipedia:", "/Portal:", "/Talk:"}
	for _, st := range ignoredLinks {
		if strings.Contains(url, st) {
			return true
		}
	}
	return false
}

// Convert judul artikel menjadi bentuk tanpa spasi
func convertToTitleCase(s string) string {
	return strings.ReplaceAll(s, " ", "_")
}

// Convert URL menjadi bentuk judul artikel
func convertToArticleTitle(URL string) string {
	s := URL[6:]
	return strings.ReplaceAll(s, "_", " ")
}
