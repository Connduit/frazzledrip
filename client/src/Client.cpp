/* TODO:
 * x86_64-w64-mingw32-gcc -nostdlib -nostartfiles -s -Wl,--entry=main -Wl,--strip-all client.c -o client.exe
 * 
 * // Working flags below, 5/72
 * x86_64-w64-mingw32-gcc -s -Wl,--entry=main -Wl,--strip-all client.c -o client.exe 
 * */

// TODO: eventually convert to c code then asm then shell code?


// #include "Client.h"
// #include "ComponentFactory.h"
// #include "TransportLayer.h"
//
// #include <iostream>
//
//
//
// Client::Client()
// {
//     // TODO: setup/configure logging
//     // loadConfig
//     // setupMessaging
//     // setupEvents
//     // setupTasks
//
//
//
//     setupComponents();
// }
//
//
//
// Client::~Client()
// {
//     if (messageHandler_)
//     {
//         delete messageHandler_;
//         messageHandler_ = 0;
//     }
//
//     if (transportLayer_)
//     {
//         delete transportLayer_;
//         transportLayer_ = 0;
//     }
//
//     [>
//     if (transportLayerFactory_)
//     {
//         delete transportLayerFactory_;
//     }*/
//
//     [>
//     delete messageHandlerPtr_;
//     messageHandlerPtr_ = 0;
//
//     delete transportLayerFactory_;
//     transportLayerFactory_ = 0;
//     */
//
// }
//
// bool Client::run()
// {
//     //Config* config = loadConfig();
//     //Config config = loadConfig();
//
//     // 1. manually resolve apis
//     [>
//     TODO:
//     would i be more effient to find the function/procedure names all at once for a single dll/module so i don't have
//     to loop through the array of functions every single time?
//
//     something like this:
//     resolve_all_functions(dll_base, hashes, pointers, count);
//
//     not sure if this is worth doing for dlls/modules tho
//
//     */
//
//     while (true)
//     {
//         if (!transportLayer_->isConnected())
//         {
//             transportLayer_->connect();
//         }
//
//         if (transportLayer_->isConnected())
//         {
//             transportLayer_->beacon();
//         }
//         Sleep(5000);
//     }
//     return false;
// }
//
//
//
// //bool Client::setupComponents(Config config)
// bool Client::setupComponents()
// {
//     messageHandler_ = new MessageHandler();
//     //transportLayerFactory_ = new TransportLayerFactory(*messageHandlerPtr_);
//     // TODO: factory needs to be passed to somewhere else? maybe the transportlayer itself?
//     //transportLayer_ = new TransportLayer(*transportLayerFactory_); // TODO:
//     //transportLayer_ = new TransportLayer(); // TransportLayer(messageHandlerPtr_);
//
//     // TODO: instead of hardcoding transport type, look at config to decide
//     transportLayer_ = TransportLayerFactory::create(TransportLayerType::TCP);
//
//
//
//     // template<class TransportLayerType> class TransportLayerFactory
//     // Arch::OwningPtr<TransportLayerType> createTransportLayer();
//     // Arch::OwningPtr<TransportLayerType> createTransportLayer(const TransportId?);
// }
//
//
// bool Client::setupTasks()
// {
//     // TODO: make an object called task? or individual object called runTask that calls the logic inside Client.run() ?
// }
//
/* let tasks (things listening for events to be risen) know
 * that subsystem is transitioning to RUNNING
 * void prepareToRun():
 *
 *
 * */
//
//
/* NOTES:
 * typedef TransportLayerFactory<TCP> TCPFactory?
 * typedef ComponentFactory<ComponentType> ComponentTypeFactory
 *
 * */
//
//
// void setupMessaging()
// {
//     messageConsumer_ = new MessageConsumer(dds_participant1, dds_participant2);
//     interalMessageConsumer_ = new InternalMessageConsumer(dds_participant1);
//     messagePublisher_ = new messagePublisher(dds_participant1, dds_participant2);
//
//     // subscriber for receiving dwell reports
//     systemResponseSubscriber = new DwellSubscriber(dds_participant1);
//
//     testActivityRequestFactory_ = new ActivityRequestFactory(dds_participant1, config_stuff);
//     testDataResponseSubscriber_ new ActivityManagerStatusSubscriber(dds_participant1, config_stuff);
// }
//
// void setupEvents()
// {
//     // allocate memory for event pool sizes
// }
//
// void setupTasks()
// {
//     scheduler_ = new TestAssistScheduler(*messagePublisher);
//
//     systemStatusTask = new TestAssistSystemStatusTask(config_stuff, *messageConsumer_);
//
//     // more task stuff...
// }
