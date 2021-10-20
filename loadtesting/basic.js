import http from 'k6/http';
import { sleep } from 'k6';

export default function() {
    var url = 'http://localhost:8000/orders'
    var payload = JSON.stringify({
        "customer": pickName(),
        "pastry": pickPastry(),
    })
    var params = {
        headers: {
            "Content-Type": "applicaion/json"
        }
    }
    http.post(url, payload, params);
    sleep(1);
}

function pickName() {
    let names = [
        "homer",
        "fry",
        "tuca",
        "bertie",
    ]

    return names[Math.floor(Math.random() * names.length)]
}

function pickPastry() {
    let pastries = [
        "croissant",
        "kouign-amann",
        "la bombe",
        "escargot",
        "almond croissant",
        "profiterole",
        "caramel tart",
    ]

    return pastries[Math.floor(Math.random() * pastries.length)]
}