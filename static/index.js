const url = "ws://" + window.location.host + "/ws";
const ws = new WebSocket(url);

const fileInput = document.querySelector("#file-input");
const fileSubmitButton = document.querySelector("#file-submit");

const driveFound = (driveName) => {
  document.querySelector("#drive-found").innerHTML = driveName;
  fileInput.disabled = false;
};

const driveNotFound = () => {
  document.querySelector("#drive-found").innerHTML = "Waiting for you...";
  fileInput.disabled = true;
};

const clear = () => {
  fileInput.value = "";
  fileInput.disabled = true;
  fileSubmitButton.disabled = true;
  fileSubmitButton.innerHTML = "Saved! You can remove your flash drive";
  setTimeout(() => {
    fileSubmitButton.innerHTML = "Save to drive";
  }, 3000);
};

fileInput.addEventListener("change", () => {
  console.log("file attached");
  console.log(fileInput.files);
  if (fileInput.files.length > 0) {
    fileSubmitButton.disabled = false;
  }
});

fileSubmitButton.addEventListener("click", async () => {
  if (fileInput.files.length > 0) {
    const file = fileInput.files[0];
    const formData = new FormData();
    formData.append("file", file);
    await fetch("/", {
      method: "POST",
      body: formData,
    });
  }

  clear();
});

ws.onmessage = function (msg) {
  let payload = {};
  try {
    payload = JSON.parse(msg.data);
  } catch (e) {
    console.log(e);
    return;
  }

  switch (payload.type) {
    case "drive":
      const driveName = payload.data;
      driveFound(driveName);
      break;
    default:
      driveNotFound();
  }
};
