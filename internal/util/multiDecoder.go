package util

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

func YamlMultiTemplate(filename string, data any) ([]*unstructured.Unstructured, error) {
	tpl, err := template.ParseFiles(filename)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}
	writer := new(bytes.Buffer)
	err = tpl.Execute(writer, data)
	if err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}
	return YamlMultiDecodeToUnstructured(writer.Bytes())
}

func YamlMultiDecodeToUnstructured(yamlBytes []byte) ([]*unstructured.Unstructured, error) {
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(yamlBytes), 1000)
	docCount := 0
	var result []*unstructured.Unstructured
	for {
		docCount++
		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil {
			if err == io.EOF {
				break
			}
			return result, fmt.Errorf("error deconding yaml document #%d: %w", docCount, err)
		}
		if rawObj.Raw == nil {
			// empty yaml doc
			continue
		}

		obj, _, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("error deconding rawObj into UnstructuredJSONScheme in document #%d: %w", docCount, err)
		}
		u, ok := obj.(*unstructured.Unstructured)
		if !ok {
			unstructuredData, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
			if err != nil {
				return nil, fmt.Errorf("error converting obj to unstructured in document #%d: %w", docCount, err)
			}

			u = &unstructured.Unstructured{Object: unstructuredData}
		}

		result = append(result, u)
	}
	return result, nil
}
