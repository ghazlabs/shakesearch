package doc_test

import (
	"strings"
	"testing"

	"pulley.com/shakesearch/internal/doc"
)

func TestGetShortHTML(t *testing.T) {
	// define test cases
	testCases := []struct {
		Name      string
		Text      string
		Query     string
		ExpResult string
	}{
		{
			Name: "Test Single Word Query",
			Text: `
				Project Gutenberg’s The Complete Works of William Shakespeare, by William Shakespeare

				This eBook is for the use of anyone anywhere in the United States and
				most other parts of the world at no cost and with almost no restrictions
				whatsoever.  You may copy it, give it away or re-use it under the terms
				of the Project Gutenberg License included with this eBook or online at
				www.gutenberg.org.  If you are not located in the United States, you’ll
				have to check the laws of the country where you are located before using
				this ebook.
				
				
				Title: The Complete Works of William Shakespeare
				
				Author: William Shakespeare
				
				Release Date: January 1994 [EBook #100]
				Last Updated: August 6, 2020
				
				Language: English
				
				Character set encoding: UTF-8
				
				*** START OF THIS PROJECT GUTENBERG EBOOK THE COMPLETE WORKS OF WILLIAM SHAKESPEARE ***
				
				
				
				
				The Complete Works of William Shakespeare
				
				
				
				by William Shakespeare
				
				
				
				
					Contents
				
				
				
							THE SONNETS
				
							ALL’S WELL THAT ENDS WELL
				
							THE TRAGEDY OF ANTONY AND CLEOPATRA
				
							AS YOU LIKE IT
				
							THE COMEDY OF ERRORS
							
			`,
			Query:     "Gutenberg",
			ExpResult: "Project <b>Gutenberg</b>’s The Complete Works of William Shakespeare, by William Shakespeare ... of the Project <b>Gutenberg</b> License included with this eBook or online at www.guten...",
		},
		{
			Name: "Test Two Words Query",
			Text: `
				THE TRAGEDY OF HAMLET, PRINCE OF DENMARK

				THE FIRST PART OF KING HENRY THE FOURTH

				THE SECOND PART OF KING HENRY THE FOURTH

				THE LIFE OF KING HENRY THE FIFTH

				THE FIRST PART OF HENRY THE SIXTH

				THE SECOND PART OF KING HENRY THE SIXTH

				THE THIRD PART OF KING HENRY THE SIXTH

				KING HENRY THE EIGHTH

				KING JOHN

				THE TRAGEDY OF JULIUS CAESAR

				THE TRAGEDY OF KING LEAR
			`,
			Query:     "King Henry",
			ExpResult: "THE FIRST PART OF <b>KING</b> <b>HENRY</b> THE FOURTH THE SECOND PART OF <b>KING</b> <b>HENRY</b> THE FOURTH THE LIFE OF <b>KING</b> <b>HENRY</b> THE FIFTH THE FIRST PART OF <b>HENRY</b> THE SIXTH THE SECOND PART OF KIN...",
		},
		{
			Name: "Test More Than One Words Query",
			Text: `
				ENOBARBUS.
				Under a compelling occasion, let women die. It were pity to cast them
				away for nothing, though, between them and a great cause they should be
				esteemed nothing. Cleopatra, catching but the least noise of this, dies
				instantly. I have seen her die twenty times upon far poorer moment. I
				do think there is mettle in death which commits some loving act upon
				her, she hath such a celerity in dying.
				
				ANTONY.
				She is cunning past man’s thought.
				
				ENOBARBUS.
				Alack, sir, no; her passions are made of nothing but the finest part of
				pure love. We cannot call her winds and waters sighs and tears; they
				are greater storms and tempests than almanacs can report. This cannot
				be cunning in her; if it be, she makes a shower of rain as well as
				Jove.
				
				ANTONY.
				Would I had never seen her!
				
				ENOBARBUS.
				O, sir, you had then left unseen a wonderful piece of work, which not
				to have been blest withal would have discredited your travel.
				
				ANTONY.
				Fulvia is dead.
				
				ENOBARBUS.
				Sir?
				
				ANTONY.
				Fulvia is dead.			
			`,
			Query:     "Cleoparta dies instantly",
			ExpResult: "esteemed nothing. <b>Cleopatra</b>, catching but the least noise of this, <b>dies</b> <b>instantly</b>. I have seen her die twenty times upon far poorer moment. I",
		},
	}
	for _, testCase := range testCases {
		// initialize document
		d, err := doc.New(
			doc.Configs{
				Lines: strings.Split(testCase.Text, "\n"),
				ShortTag: doc.Tag{
					Start: "<b>",
					End:   "</b>",
				},
			},
		)
		if err != nil {
			t.Fatalf("unable to initialize new document due: %v", err)
		}
		// get short html
		result := d.GetShortHTML(testCase.Query)
		if result != testCase.ExpResult {
			t.Fatalf("unexpected result, exp: %v, got: %v", testCase.ExpResult, result)
		}
	}
}

func TestGetHighlightedHTML(t *testing.T) {
	// define test cases
	testCases := []struct {
		Name      string
		Text      string
		Query     string
		ExpResult string
	}{
		{
			Name:      "Test Any Case Highlighted",
			Text:      "King kING KING kIng",
			Query:     "King",
			ExpResult: `<span style="highlight">King</span> <span style="highlight">kING</span> <span style="highlight">KING</span> <span style="highlight">kIng</span>`,
		},
		{
			Name: "Test Substring Highlighted",
			Text: `
				Project Gutenberg’s The Complete Works of William Shakespeare, by William Shakespeare

				This eBook is for the use of anyone anywhere in the United States and
				most other parts of the world at no cost and with almost no restrictions
				whatsoever.  You may copy it, give it away or re-use it under the terms
				of the Project Gutenberg License included with this eBook or online at
				www.gutenberg.org.  If you are not located in the United States, you’ll
				have to check the laws of the country where you are located before using
				this ebook.
			`,
			Query: "Gutenberg",
			ExpResult: `
				Project <span style="highlight">Gutenberg</span>’s The Complete Works of William Shakespeare, by William Shakespeare

				This eBook is for the use of anyone anywhere in the United States and
				most other parts of the world at no cost and with almost no restrictions
				whatsoever.  You may copy it, give it away or re-use it under the terms
				of the Project <span style="highlight">Gutenberg</span> License included with this eBook or online at
				www.<span style="highlight">gutenberg</span>.org.  If you are not located in the United States, you’ll
				have to check the laws of the country where you are located before using
				this ebook.
			`,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// initialize document
			d, err := doc.New(doc.Configs{})
			if err != nil {
				t.Fatalf("unable to initialize document due: %v", err)
			}
			// get highlighted html
			result := d.GetHighlightedHTML(testCase.Query)
			if result != testCase.ExpResult {
				t.Fatalf("unexpected result, exp: %v, got: %v", testCase.ExpResult, result)
			}
		})
	}
}
