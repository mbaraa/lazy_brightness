let b;

const client = new XMLHttpRequest();
client.open("GET", "/brits/get", false);
client.send();

b = JSON.parse(client.responseText);
document.getElementById("brightness").value = (100 * b["brightness"]) / b["max"]

console.log("b", b);

let devices;

client.open("GET", "/brits/devices", false);
client.send();

devices = JSON.parse(client.responseText)["devices"];
console.log("devices", devices);

function putDevices() {
	let devicesSelections = "";

	for (device of devices) {
		devicesSelections += `<option>${device}</option>`;
	}
	
	document.getElementById("devices").innerHTML = devicesSelections;
}
putDevices()

function setDevice() {
	fetch(`/brits/set_device?device=${document.getElementById("devices").value}`, {
		method: "GET",
	})
}

function setBrightness() {
	fetch(`/brits/set?b=${document.getElementById("brightness").value}`, {
		method: "GET",
	})
}

function decBrightness() {
	fetch("/brits/dec", {
		method: "GET"
	})
}

function incBrightness() {
	fetch("/brits/inc", {
		method: "GET"
	})
}