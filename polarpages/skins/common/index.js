function getCookie(cname) {
  let name = cname + "=";
  let decodedCookie = decodeURIComponent(document.cookie);
  let ca = decodedCookie.split(';');
  for(let i = 0; i <ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

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
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then(response => {
      if (response.status === 404) {
        return createPage(data)
      } else if (!response.ok) {
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

const createPage = (data) => {
  fetch("/api/e/" + data.Title, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
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

const ggCreateAccount = () => {
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;
  const confirmPassword = document.getElementById('confirm-password').value;

  if (username === '' || password === '') {
    return;
  }

  // Check if password is confirmed
  if (password !== confirmPassword) {
    return;
  }
  
  const requestBody = {
    username: username,
    password: password
  };

  fetch('/api/CreateAccount', {
    method: 'POST',
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(requestBody)
  })
    .then(response => {
      if (!response.ok) {
        console.log(response.json())
        throw new Error("Response not ok");
      }
      return response.json();
    })
    .then(data => {
      console.log("Succesfully created user account " + username);
      console.log(data);

      // Redirect to main page if succesful
      window.location.href = '/wiki/main_page';
    })
    .catch(error => {
      console.error(error)
    })
}

const ggLogin = () => {
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;

  if (username === '' || password === '') {
    return;
  }

  const requestBody = {
    username: username,
    password: password
  };

  fetch('/api/Login', {
    method: 'POST',
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(requestBody)
  })
    .then(response => {
      if (!response.ok) {
        throw new Error("Response not ok");
      }
      return response.json();
    })
    .then(data => {
      console.log("Succesfully logged in as user " + username);
      console.log(data);

      // Redirect to main page if succesful
      window.location.href = '/wiki/main_page';
    })
    .catch(error => {
      console.error(error)
    })
}

