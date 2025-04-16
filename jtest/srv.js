import * as io from "io2";
import * as os from "os";
import * as http from "http";

const cli = new http.Client()
const resp = cli.get("https://baidu.com")
console.log(resp.statusCode)

let out;
try {
    out = os.create("out.html", 0o755)
    const n = io.copy(out, resp.body)
    console.log(">>>> ", n)
} finally {
    if (out) {
        out.close()
    }
}