const toggleView = () => {
  document.getElementById("gg-page-nav").style.display = document.getElementById("gg-page-nav").style.display === "none" ? "block" : "none";
};

const sendEditRequest = () => {
  let title = window.location.pathname.split('/').pop();
  let content = document.querySelector('.gg-page-source').value;

  let data = {
    Title: title,
    Content: content
  };

  fetch("/api/e/" + title, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(data),
  })
    .then(response => {
      if (!response.ok) {
        throw new Error("Response not ok");
      }
      return response.json();
    })
    .then(data => {
      console.log(data)
    })
    .catch(error => {
      console.error(error)
    })
}
