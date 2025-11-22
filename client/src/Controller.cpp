#include "Controller.h"

#include <iostream>


void Controller::handleDefault(const InternalMessage& msg)
{
	std::cout << "Controller::handleDefault()" << std::endl;
}

void Controller::handleNone(const InternalMessage& msg)
{
	std::cout << "Controller::handleNone()" << std::endl;
}

void Controller::handleExecuteCommand(const InternalMessage& msg)
{
	std::cout << "Controller::handleExecuteCommand()" << std::endl;
	// std::cout << "executeCommand: " << byte2string(data).c_str() << std::endl;


	// TODO: instead add result into the queue, and have transporter send them all when ready?
	// TODO: if the original (identified by the messageId) incoming internalMessage was marked as high priority, 
	//		 then skip the queue and send the message immediately 

	//std::string command = byte2string(data).c_str();
	//std::string result;

	/*
	// EASIEST METHOD - just use _popen directly
	FILE* pipe = _popen(command.c_str(), "r");
	if (!pipe)
	{
		result = "ERROR: Failed to execute command";
	}
	else
	{
		char buffer[128];
		while (fgets(buffer, sizeof(buffer), pipe) != NULL)
		{
			result += buffer;
		}
		_pclose(pipe);
	}

	//queueResponse(COMMAND_RESULT, result, msg_id); // TODO: need to pass all of InternalMessage if i want ID
	std::cout << result << std::endl;

	// TODO: props should make a separate function for InternalMessage generation
	InternalMessage outMsg;
	outMsg.data_ = string2byte(result);

	MessageHeader header;
	header.messageType_ = MessageType::COMMAND_RESULT;
	header.messageId_ = 105; // TODO: 
	header.dataSize_ = outMsg.data_.size();
	outMsg.header_ = header;

	// TODO: 
	// bool sendResult = transportLayer_->sendMessage(outMsg);
	// return sendResult;
	std::cout << "TODO: " << std::endl;
	return false;
	*/

	// TODO: add callback function here for subsystem to register handlers 
	// for sending outgoing messages? 
}
