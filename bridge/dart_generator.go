package custom_libp2p_bridge

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
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
  static final List<Function(String)> _onErrorCallbacks = [defaultErrorCallback];

  %vException(this.cause) {
    for (var callback in _onErrorCallbacks) {
	  callback(cause);
    }
  }

  static void defaultErrorCallback(String error) {	
    print("%vException: $error");
  }

  static void addOnErrorCallback(Function(String) callback) {
    _onErrorCallbacks.add(callback);
  }

  static void removeOnErrorCallback(Function(String) callback) {
    _onErrorCallbacks.remove(callback);
  }
}

class %v {
  static const MethodChannel _channel = const MethodChannel('custom_libp2p');
  static bool bindingEnabled = Binding().isSupported();
  static final List<Function(models.Warning)> _onWarningCallbacks = [];
  static final List<Function(String)> _onInfoCallbacks = [];
  static final List<Function(models.Message)> _onMessageCallbacks = [];
  static bool _notificationHandlerRunning = true;

  static Future<Uint8List> _handleRawResult(List<Uint8List> result) async {
    switch (result.length) {
      case 1:
        return result[0];
      case 2:
        var error = new String.fromCharCodes(result[1]);
        throw new CustomLibP2PException(error);
      default:
        throw new CustomLibP2PException("Unexpected result: $result");
    }
  }

  static Future<Uint8List> _call(String name, Uint8List payload) async {
    if (bindingEnabled) {
      List<Uint8List> result = await Binding().callAsync(name, payload);
      return _handleRawResult(result);
    }
    List<Uint8List> result = await _channel.invokeMethod(name, payload);
    return _handleRawResult(result);
  }

  static void startNotificationHandler() async {
    while (_notificationHandlerRunning) {
      var notifications = await _flushNotifications();
      for (var notification in notifications.notification) {
        if (notification.hasWarningNotification()) {
          _onWarningCallbacks.forEach(
              (callback) => callback(notification.warningNotification));
        } else if (notification.hasInfoNotification()) {
          _onInfoCallbacks
              .forEach((callback) => callback(notification.infoNotification));
        } else if (notification.hasMessageNotification()) {
          _onMessageCallbacks.forEach(
              (callback) => callback(notification.messageNotification));
        } else {
          print("Unhandled notification: $notification");
        }
      }
    }
  }

  static void stopNotificationHandler() {
    _notificationHandlerRunning = false;
  }

  static void addOnWarningCallback(Function(models.Warning) callback) {
    _onWarningCallbacks.add(callback);
  }

  static void addOnInfoCallback(Function(String) callback) {
    _onInfoCallbacks.add(callback);
  }

  static void addOnMessageCallback(Function(models.Message) callback) {
    _onMessageCallbacks.add(callback);
  }

  static void removeOnWarningCallback(Function(models.Warning) callback) {
    _onWarningCallbacks.remove(callback);
  }

  static void removeOnInfoCallback(Function(String) callback) {
    _onInfoCallbacks.remove(callback);
  }

  static void removeOnMessageCallback(Function(models.Message) callback) {
    _onMessageCallbacks.remove(callback);
  }
`, DartBridgeClassName, DartBridgeClassName, DartBridgeClassName, DartBridgeClassName)
	keys := make([]string, 0)
	for k := range bridgeMapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, name := range keys {
		info := bridgeMapping[name]
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
	code += "}\n"
	return os.WriteFile(outputPath, []byte(code), 0644)
}

func GenerateConstants(outputPath string) error {
	code := `// Auto Generated Code. Do Not Edit.`
	for varName, err := range exportedErrors {
		code += fmt.Sprintf("\nvar %s = '%s';", varName, err)
	}
	return os.WriteFile(outputPath, []byte(code), 0644)
}
