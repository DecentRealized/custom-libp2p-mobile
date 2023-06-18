import * as flatbuffers from 'flatbuffers';
import {GetHelloMessageRequest} from "../libs/model/get-hello-message-request";
import {StringResponse} from "../libs/model/string-response";

export const HelloWorldSample = async () => {
    const builder = new flatbuffers.Builder(0);
    const userName = builder.createString('Wasm Tester');
    GetHelloMessageRequest.startGetHelloMessageRequest(builder);
    GetHelloMessageRequest.addUserName(builder, userName);
    const offset = GetHelloMessageRequest.endGetHelloMessageRequest(builder);
    builder.finish(offset);

    const bytes = builder.asUint8Array()

    console.log('request', bytes);
    const rawResponse = await sendToWorker('getHelloWorld', bytes)
    console.log(rawResponse)
    const responseBuffer = new flatbuffers.ByteBuffer(rawResponse);
    const response = StringResponse.getRootAsStringResponse(responseBuffer);
    // if (response.error()) {
    //     throw new Error(response.error()!)
    // }
    return response
}

let counter = 0;
const sendToWorker = (name: string, request: Uint8Array) => {
    const myWorker = new Worker('worker.js');
    counter++;
    const id = counter.toString()

    return new Promise<Uint8Array>((resolve, reject) => {

        const callbackError = (e: any) => {
            reject('callbackError: ' + e.message)
        }
        const callbackMessageError = (e: any) => {
            reject('callbackMessageError: ' + e)
        }
        const callback = (e: any) => {
            const data = e.data || {}
            if (id !== data.id) {
                // if not same if we should not reject
                return
            }
            myWorker.removeEventListener('message', callback)
            const {error, response} = data;
            if (error) {
                reject(error)
            }
            resolve(response);
        }

        myWorker.addEventListener('message', callback)
        myWorker.addEventListener('error', callbackError)
        myWorker.addEventListener("messageerror", callbackMessageError)
        myWorker.postMessage({id, name, request});
    })
}