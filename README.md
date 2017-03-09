## envoke

A demo client-side application for persisting music metadata and rights to BigchainDB/IPDB.

### Install 

Download and install [Go](https://golang.org/dl/).

In a terminal window, `go get github.com/zbo14/envoke/cmd/envoke`

### Usage

In a terminal window, `cd ~/go/src/github.com/zbo14/envoke`...

* **Demo app**
	
	`sh start_app.sh` 

	You will be prompted to enter an endpoint to the BigchainDB/IPDB http-api. 

	In your browser, go to `http://localhost:8888/<endpoint>`
    
    Endpoints:  		
    - login_register,
    - compose_publish,
    - record_release,
    - right_license

*  **Run tests**

	`sh run_tests.sh`

	You will be prompted to enter an endpoint to the BigchainDB/IPDB http-api and a path to an audio file.