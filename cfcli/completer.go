package cfcli

import (
	"reflect"
	"regexp"
	"strings"

	"code.cloudfoundry.org/cli/command/common"
	"github.com/c-bata/go-prompt"
)

func createCommandCompletions() []prompt.Suggest {
	t := reflect.TypeOf(common.Commands)
	completions := make([]prompt.Suggest, 0, t.NumField()+len(cfShellCmds))
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		command := f.Tag.Get("command")
		alias := f.Tag.Get("alias")
		description := f.Tag.Get("description")

		if command == "" {
			continue
		}
		if strings.HasPrefix(command, "v3-") {
			continue
		}
		completions = append(completions, prompt.Suggest{Text: command, Description: description})

		if alias != "" {
			completions = append(completions, prompt.Suggest{Text: alias, Description: description})
		}
	}
	for _, cmd := range cfShellCmds {
		completions = append(completions, prompt.Suggest{Text: cmd[0], Description: cmd[1]})
	}
	return completions
}

func createUsageCompletionsMap() (map[string]string, map[string]string) {
	usageMap := make(map[string]string)
	aliasMap := make(map[string]string)
	t := reflect.TypeOf(common.Commands)
	for i := 1; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type.Kind() == reflect.Struct {
			usageVar, exists := f.Type.FieldByName("usage")
			if !exists {
				continue
			}
			command := f.Tag.Get("command")
			alias := f.Tag.Get("alias")
			usage := usageVar.Tag.Get("usage")
			usageMap[command] = usage
			aliasMap[alias] = command
		}
	}
	return usageMap, aliasMap
}

func createUsageCompletion(commandText string) []prompt.Suggest {
	commandWords := strings.Split(commandText, " ")
	cmd := commandWords[0]

	usageMap, aliasMap := createUsageCompletionsMap()
	command, exists := aliasMap[cmd]
	if !exists {
		command = cmd
	}
	usage := usageMap[command]

	var commandUsage string
	for _, line := range strings.Split(usage, "\n") {
		if strings.Contains(line, "CF_NAME") || strings.HasPrefix(line, "cf ") {
			commandUsage = strings.Replace(line, "CF_NAME", "", -1)
			commandUsage = strings.TrimLeft(commandUsage, "cf")
			commandUsage = strings.Replace(commandUsage, command, "", 1)
			commandUsage = strings.TrimSpace(commandUsage)
			break
		}
	}
	commandArgs := strings.Split(commandUsage, " ")

	var currentArg string
	if len(commandWords)-2 < len(commandArgs) {
		currentArg = commandArgs[len(commandWords)-2]
	}

	if IsResolvableKeyWord(currentArg) {
		entities := ResolveKeyWord(currentArg)

		filteredEntities := make([]string, 0)
		lastCommandWord := commandWords[len(commandWords)-1]

		for _, entity := range entities {
			if lastCommandWord == "" || strings.Contains(strings.ToUpper(entity), strings.ToUpper(lastCommandWord)) {
				filteredEntities = append(filteredEntities, entity)
			}
		}

		suggestions := []prompt.Suggest{}
		for _, entity := range filteredEntities {
			suggestions = append(suggestions, prompt.Suggest{Text: entity})
		}
		return suggestions
	}

	var remainingArgs string
	if len(commandWords)-2 < len(commandArgs) {
		remainingArgs = strings.Join(commandArgs[len(commandWords)-2:], " ")
	}

	if remainingArgs == "" {
		return []prompt.Suggest{}
	} else {
		return []prompt.Suggest{{Text: " ", Description: remainingArgs}}
	}
}

func Completer(d prompt.Document) []prompt.Suggest {
	cleanText := removeDupSpaces(d.Text)

	if cleanText == "" {
		return []prompt.Suggest{}
	} else {
		words := strings.Split(cleanText, " ")
		if len(words) > 1 {
			return createUsageCompletion(cleanText)
		} else {
			suggestions := createCommandCompletions()
			suggestions = moveExactMatchToFront(suggestions, words[0])
			return prompt.FilterContains(suggestions, words[0], true)
		}
	}
}

func moveExactMatchToFront(suggestions []prompt.Suggest, keyword string) []prompt.Suggest {
	var result []prompt.Suggest

	for i, suggestion := range suggestions {
		if suggestion.Text == keyword {
			result = append(result, suggestion)
			if i > 0 {
				// append everything before i
				result = append(result, suggestions[0:i-1]...)
			}
			if i < len(suggestions)-1 {
				// append everything after i
				result = append(result, suggestions[i+1:]...)
			}
			break
		}
	}

	if len(result) == 0 {
		result = append(result, suggestions...)
	}

	return result
}

func removeDupSpaces(input string) string {
	regTrialingSpaces := regexp.MustCompile(`[\s\p{Zs}]{2, }$`)
	regMiddleSpaces := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	final := strings.TrimLeft(input, " ")
	final = regTrialingSpaces.ReplaceAllString(final, " ")
	final = regMiddleSpaces.ReplaceAllString(final, " ")

	return final
}
