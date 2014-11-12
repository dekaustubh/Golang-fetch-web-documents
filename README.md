Golang-fetch-web-documents
==========================

A simple golang code to fetch web documents

Problem Statement : A program that fetches 20 URLs 3 at a
time and computes the average byte size of all the fetched documents

Solution : The solution for the problem definition is Fetcher package which provides Fetcher object that accepts list of urls and contains members that indicate the results of fetching at the end of the operation.
Now Fetcher is designed considering the concurrency into the picture. As functionality involves network operation and all the three job executors are fetching the urls from a common source, there is requirement of synchronisation. For this read URL operation we have used Mutex provided by sync package of the standard go library.  
To communicate between our main routine and spawned go routines we are using channels. 
To keep three jobs executing always we have used select operator over channels that are passed to the go routines.
