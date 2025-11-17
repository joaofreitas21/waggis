document.addEventListener("DOMContentLoaded", function () {
  let isEmailFormOpen = false;


  // Handle click on "Email me!" button
  document.body.addEventListener("click", async function (e) {
    const emailBtn = e.target.closest("[data-action='email']");
    if (emailBtn) {
      e.preventDefault();

      if (isEmailFormOpen) {
        closeEmailForm();
        isEmailFormOpen = false;
        return;
      }

      await showEmailForm();
      isEmailFormOpen = true;
    }
  });
});

async function showEmailForm() {
  const slot = document.getElementById("email-form-slot");
  const globe = document.getElementById("globe-wrap");
  const emailBtn = document.querySelector("[data-action='email']");

  if (!slot) return;

  // Get color from email button or default to purple
  const colorClass = emailBtn?.getAttribute("data-color") || "text-purple-400";

  // Apply color to slot
  slot.classList.forEach((c) => {
    if (c.startsWith("text-")) slot.classList.remove(c);
  });
  slot.classList.add(colorClass);

  // Fetch and render form
  const resp = await fetch("/email-form");
  if (resp.ok) {
    slot.innerHTML = await resp.text();
    slot.classList.remove("hidden");
  } else {
    slot.innerHTML =
      "<div class='text-red-400 p-6'>Failed to load form :(</div>";
    slot.classList.remove("hidden");
  }

  document.documentElement.classList.add("email-open");
  document.body.classList.add("email-open");

  // Hide globe, keep info-line visible
  if (globe) globe.classList.add("hidden");
  document.body.style.overflow = "hidden";
}

function closeEmailForm() {
  const slot = document.getElementById("email-form-slot");
  const globe = document.getElementById("globe-wrap");
  const form = document.getElementById("email-contact-form");

  if (slot) {
    slot.innerHTML = "";
    slot.classList.add("hidden");
  }

  // Reset form state
  if (form) {
    form.reset();
    closeEmailMessage();
  }

  document.documentElement.classList.remove("email-open");
  document.body.classList.remove("email-open");

  // Show globe, restore scroll
  if (globe) globe.classList.remove("hidden");
  document.body.style.overflow = "";
}

async function submitEmailForm(event) {
  event.preventDefault();

  const form = document.getElementById("email-contact-form");
  const submitBtn = document.getElementById("email-submit-btn");
  
  const messageSlot = document.getElementById("email-message-slot");

  if (!form || !submitBtn || !messageSlot) return;

  
  closeEmailMessage();

  // Disable submit button and show loading state
  submitBtn.disabled = true;
  submitBtn.textContent = "Sending...";

  
  const formData = {
    name: document.getElementById("email-name").value.trim(),
    email: document.getElementById("email-email").value.trim(),
    subject: document.getElementById("email-subject").value.trim(),
    message: document.getElementById("email-message").value.trim(),
  };

  try {
    const resp = await fetch("/api/email", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
    });

    const data = await resp.json();

    if (data.success) {
      // Fetch and show success card
      const message = data.message || "Thanks for reaching out, I will be in contact soon!";
      const successResp = await fetch(
        `/email-form/success?message=${encodeURIComponent(message)}`
      );
      if (successResp.ok) {
        messageSlot.innerHTML = await successResp.text();
      }

      // Reset form and button
      form.reset();
      submitBtn.disabled = false;
      submitBtn.textContent = "Send";
    } else {
      // Fetch and show error card
      let errorMessage = data.error || "Failed to send email. Please try again.";
      let retryAfter = data.retry_after || 0;

      const errorResp = await fetch(
        `/email-form/error?message=${encodeURIComponent(errorMessage)}&retry_after=${retryAfter}`
      );
      if (errorResp.ok) {
        messageSlot.innerHTML = await errorResp.text();
      }

      // Re-enable submit button
      submitBtn.disabled = false;
      submitBtn.textContent = "Send";
    }
  } catch (error) {
    // Network or other error
    const errorMessage = "Network error. Please check your connection and try again.";
    const errorResp = await fetch(
      `/email-form/error?message=${encodeURIComponent(errorMessage)}&retry_after=0`
    );
    if (errorResp.ok) {
      messageSlot.innerHTML = await errorResp.text();
    }

    submitBtn.disabled = false;
    submitBtn.textContent = "Send";
  }
}

function closeEmailMessage() {
  const messageSlot = document.getElementById("email-message-slot");
  if (messageSlot) {
    messageSlot.innerHTML = "";
  }
}


window.closeEmailForm = closeEmailForm;
window.submitEmailForm = submitEmailForm;
window.closeEmailMessage = closeEmailMessage;