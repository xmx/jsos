import * as http from "http";

const mux = http.newServeMux();
mux.handleFunc("/ping", (w, r) => {
    console.log(r.remoteAddr)
    w.header().set("Content-Type", "application/json")
    w.writeHeader(http.statusAccepted)

    const data = {
        name: "PING",
        date: new Date(),
    }

    w.write(JSON.stringify(data))
})

http.listenAndServe(":9099", mux)
