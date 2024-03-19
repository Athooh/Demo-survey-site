document.getElementById("newsletter-form").addEventListener("submit", function(event) {
  event.preventDefault();
  var formData = new FormData(this);
  var xhr = new XMLHttpRequest();
  xhr.open("POST", "/subscribe");
  xhr.onreadystatechange = function() {
      if (xhr.readyState === XMLHttpRequest.DONE) {
          if (xhr.status === 200) {
              var response = JSON.parse(xhr.responseText);
              if (response.success) {
                  alert(response.message);
                  document.getElementById("newsletter-form").reset(); // Reset form fields
              } else {
                  alert("Error: " + response.message);
              }
          } else {
              alert("Error: Unable to subscribe at the moment. Please try again later.");
          }
      }
  };
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.send(JSON.stringify(Object.fromEntries(formData.entries())));
});