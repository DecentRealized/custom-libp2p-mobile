package custom_libp2p_bridge

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

func CreateModelExporter(protocolDir string, outputPath string) error {
	protoFileExtension := ".proto"
	dartGenFileSubExtensions := []string{".pb", ".pbenum", ".pbjson", ".pbserver"}
	dartGenFileExtension := ".dart"

	var exportNames []string

	files, err := ioutil.ReadDir(protocolDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		fName := file.Name()
		if len(fName) > len(protoFileExtension) &&
			fName[len(fName)-len(protoFileExtension):] == ".proto" {
			for _, subExt := range dartGenFileSubExtensions {
				exportName := fmt.Sprintf("%s%s%s",
					fName[:len(fName)-len(protoFileExtension)],
					subExt,
					dartGenFileExtension)
				exportNames = append(exportNames, exportName)
			}
		}
	}
	code := `// Auto Generated Code. Do Not Edit.
`
	for _, exportName := range exportNames {
		code += fmt.Sprintf("\nexport '%s';", exportName)
	}
	return os.WriteFile(outputPath, []byte(code), 0644)
}

const DartBridgeClassName = "CustomLibP2P"

func CreateDartBinding(outputPath string) error {
	code := fmt.Sprintf(`// Auto Generated Code. Do Not Edit.
import 'dart:async';
import 'dart:io';

import 'package:custom_libp2p/bridge/binding_stub.dart'
    if (dart.library.io) 'package:custom_libp2p/bridge/binding.dart'
    if (dart.library.js) 'package:custom_libp2p/bridge/binding_stub.dart';
import 'package:custom_libp2p/models/models.dart' as models;
import 'package:flutter/services.dart';
import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';

class %vException implements Exception {
  String cause;

  %vException(this.cause) {
    print("%vException: $cause");
  }
}

class %v {
  static const MethodChannel _channel = const MethodChannel('custom_libp2p');
  static bool bindingEnabled = Binding().isSupported();

  static Future<Uint8List> _call(String name, Uint8List payload) async {
    if (bindingEnabled) {
      return await Binding().callAsync(name, payload);
    }
    return await _channel.invokeMethod(name, payload);
  }
`, DartBridgeClassName, DartBridgeClassName, DartBridgeClassName, DartBridgeClassName)
	for name, info := range bridgeMapping {
		if info.output != nil {
			output := reflect.TypeOf(info.output).String()
			if output[0] == '*' {
				output = output[1:]
			}
			code += fmt.Sprintf("\n  static Future<%s>", output)
		} else {
			code += "\n  static Future<void>"
		}
		code += fmt.Sprintf(" %s", name)
		if info.input != nil {
			input := reflect.TypeOf(info.input).String()
			if input[0] == '*' {
				input = input[1:]
			}
			code += fmt.Sprintf("(\n      %s input)", input)
		} else {
			code += "()"
		}
		code += " async {\n"
		code += info.preCallWrapper
		if info.input != nil {
			code += fmt.Sprintf("    var bytesInput = input.writeToBuffer();\n")
		} else {
			code += "    var bytesInput = Uint8List.fromList([]);\n"
		}
		if info.output != nil {
			code += fmt.Sprintf("    var bytesOutput = await _call('%s', bytesInput);\n", name)
			output := reflect.TypeOf(info.output).String()
			if output[0] == '*' {
				output = output[1:]
			}
			code += fmt.Sprintf("    var output = %s.fromBuffer(bytesOutput);\n", output)
			code += fmt.Sprintf("    return output;\n")
		} else {
			code += fmt.Sprintf("    await _call('%s', bytesInput);\n", name)
		}
		code += info.postCallWrapper
		code += "  }\n"
	}
	code += "}"
	return os.WriteFile(outputPath, []byte(code), 0644)
}
