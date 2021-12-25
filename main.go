package main

import (
	"log"
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
	"gopkg.in/yaml.v2"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			err := generateYamlFile(gen, f)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
}

type Messages []*Message

type Message struct {
	Name   string   `yaml:"name"`
	Fields []*Field `yaml:"fields"`
}

type Field struct {
	Name   string `yaml:"name"`
	Number int    `yaml:"number"`
}

type Services []*Service

type Service struct {
	Name    string    `yaml:"name"`
	Methods []*Method `yaml:"methods"`
}

type Method struct {
	Name       string `yaml:"name"`
	InputType  string `yaml:"input_type"`
	OutputType string `yaml:"output_type"`
}

type YamlOutput struct {
	Messages Messages `yaml:"messages"`
	Services Services `yaml:"services"`
}

func generateYamlFile(gen *protogen.Plugin, file *protogen.File) error {
	var yamlF YamlOutput
	filename := file.GeneratedFilenamePrefix + ".yaml"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	messages := getMessages(file.Messages)
	services := getServices(file.Services)

	yamlF.Messages = messages
	yamlF.Services = services
	yamlData, err := yaml.Marshal(&yamlF)
	if err != nil {
		return err
	}
	g.Write(yamlData)
	return nil
}

func getFields(m *protogen.Message) []*Field {
	var fields []*Field
	for _, field := range m.Fields {
		name := string(field.Desc.Name())
		number := int(field.Desc.Number())
		fields = append(fields, &Field{Name: name, Number: number})
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Number < fields[j].Number
	})
	return fields
}

func getMethods(s *protogen.Service) []*Method {
	var methods []*Method
	for _, method := range s.Methods {
		name := string(method.Desc.FullName())
		inputType := string(method.Desc.Input().FullName())
		outputType := string(method.Desc.Output().FullName())
		m := &Method{Name: name, InputType: inputType, OutputType: outputType}
		methods = append(methods, m)
	}
	return methods
}

func getServices(protoServices []*protogen.Service) []*Service {
	var services []*Service
	for _, service := range protoServices {
		name := string(service.Desc.FullName())
		methods := getMethods(service)
		m := Service{Name: name, Methods: methods}
		services = append(services, &m)
	}
	return services
}

func getMessages(protoMessages []*protogen.Message) []*Message {
	var messages []*Message
	for _, message := range protoMessages {
		stack := []*protogen.Message{}
		stack = append(stack, message)
		for len(stack) > 0 {
			n := len(stack) - 1
			val := stack[n]
			stack = stack[:n]
			name := string(val.Desc.FullName())
			fields := getFields(val)
			m := Message{Name: name, Fields: fields}
			messages = append(messages, &m)
			stack = append(stack, val.Messages...)
		}
	}
	return messages
}
