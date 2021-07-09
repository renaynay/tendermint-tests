FROM golang:latest

RUN git clone -b dummy https://github.com/renaynay/lazyledger-core.git
RUN cd lazyledger-core/abci/example/dummyapp && rm go.mod && go get && go build
RUN cd lazyledger-core/abci/example/dummyapp && ./dummyapp init && ./dummyapp start --dummy.txs 5 --dummy.tx-size 10
