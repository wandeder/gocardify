package anki_bot

import (
	"errors"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"os"
	"regexp"
	"time"
)

func parseMsg(msg string) (front string, back string, err error) {
	re := regexp.MustCompile(`(?i)front:\s*(.+?)\s*back:\s*(.+)`)

	matches := re.FindStringSubmatch(msg)

	if len(matches) != 3 {
		return "", "", fmt.Errorf("Incorrect format")
	}

	return matches[1], matches[2], nil
}

func CreateCard(msg string) error {
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	chromeCaps := chrome.Capabilities{
		Args: []string{"--headless=new", "--disable-gpu", "--no-sandbox"},
	}
	caps.AddChrome(chromeCaps)

	wd, err := selenium.NewRemote(caps, os.Getenv("SELENIUM_SERVER"))
	if err != nil {
		fmt.Println("Ошибка при создании удаленного веб-драйвера:", err)
		return err
	}
	defer wd.Quit()

	if err := wd.Get(os.Getenv("ANKI_URL")); err != nil {
		fmt.Println("Ошибка при открытии страницы:", err)
		return err
	}

	time.Sleep(5 * time.Second)

	emailField, err := wd.FindElement(selenium.ByXPATH, "//*[@id=\"H3pLrX8QHMmpf9mBwr3y6\"]")
	if err != nil {
		fmt.Println("Ошибка при поиске поля Email:", err)
		return err
	}
	err = emailField.SendKeys(os.Getenv("ANKI_LOGIN"))
	if err != nil {
		fmt.Println("Ошибка при вставке значения в поле Email:", err)
		return err
	}
	passwordField, err := wd.FindElement(selenium.ByXPATH, "//*[@id=\"T3mcXbGxUnhRcYXkT4irD\"]")
	if err != nil {
		fmt.Println("Ошибка при поиске поля Password:", err)
		return err
	}
	err = passwordField.SendKeys(os.Getenv("ANKI_PASSWORD"))
	if err != nil {
		fmt.Println("Ошибка при вставке значения в поле Password:", err)
		return err
	}

	err = passwordField.SendKeys(selenium.EnterKey)
	if err != nil {
		fmt.Println("Ошибка при нажатии клавиши Enter:", err)
		return err
	}

	fmt.Println("Вход выполнен успешно!")

	time.Sleep(5 * time.Second)

	currentURL, err := wd.CurrentURL()
	if err != nil {
		fmt.Println("Ошибка при получении текущего URL:", err)
		return err
	}
	if currentURL != os.Getenv("ANKI_DESK_URL") {
		return errors.New("Unsucsessful enter")
	}

	if err := wd.Get(os.Getenv("ANKI_ADD_URL")); err != nil {
		fmt.Println("Ошибка при открытии страницы:", err)
		return err
	}

	frontElement, err := wd.FindElement(selenium.ByXPATH, "/html/body/div/main/form/div[1]/div/div")
	if err != nil {
		fmt.Println("Ошибка при поиске элемента Front:", err)
		return err
	}
	front, back, err := parseMsg(msg)
	if err != nil {
		fmt.Println("Ошибка парсинга сообщения")
		return err
	}
	err = frontElement.SendKeys(front)
	if err != nil {
		fmt.Println("Ошибка при вставке значения в поле Front:", err)
		return err
	}
	backElement, err := wd.FindElement(selenium.ByXPATH, "/html/body/div/main/form/div[1]/div/div")
	if err != nil {
		fmt.Println("Ошибка при поиске элемента Back:", err)
		return err
	}
	err = backElement.SendKeys(back)
	if err != nil {
		fmt.Println("Ошибка при вставке значения в поле Back:", err)
		return err
	}
	err = passwordField.SendKeys(selenium.EnterKey)
	if err != nil {
		fmt.Println("Ошибка при нажатии клавиши Enter:", err)
		return err
	}

	fmt.Println("Card created успешно!")
	return nil
}
