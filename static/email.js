document.addEventListener("DOMContentLoaded", function () {
    let isEmailFormOpen = false;
  
    // Handle click on "Email me!" button
    document.body.addEventListener("click", async function (e) {
      const emailBtn = e.target.closest("[data-action='email']");
      if (emailBtn) {
        e.preventDefault();
  
        if (isEmailFormOpen) {
          closeEmailForm();
          isEmailFormOpen = false; // Ensure flag is reset
          return;
        }
  
        await showEmailForm();
        isEmailFormOpen = true; // Ensure flag is set
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
      hideMessages();
    }
  
    // Show globe, restore scroll
    if (globe) globe.classList.remove("hidden");
    document.body.style.overflow = "";
  }
  
  async function submitEmailForm(event) {
    event.preventDefault();
  
    const form = document.getElementById("email-contact-form");
    const submitBtn = document.getElementById("email-submit-btn");
    const successMsg = document.getElementById("email-success-message");
    const errorMsg = document.getElementById("email-error-message");
    const errorText = document.getElementById("email-error-text");
  
    if (!form || !submitBtn) return;
  
    // Hide previous messages
    hideMessages();
  
    // Disable submit button and show loading state
    submitBtn.disabled = true;
    submitBtn.textContent = "Sending...";
  
    // Collect form data
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
        // Show success message
        if (successMsg) {
          successMsg.classList.remove("hidden");
        }
  
        // Reset form
        form.reset();
  
        // Auto-close after 10 seconds
        setTimeout(() => {
          closeEmailForm();
          
          if (window.isEmailFormOpen !== undefined) {
            window.isEmailFormOpen = false;
          }
        }, 10000);
      } else {
        // Show error message
        if (errorMsg && errorText) {
          let errorMessage =
            data.error || "Failed to send email. Please try again.";
  
          // Add retry time if rate limited
          if (data.retry_after && data.retry_after > 0) {
            errorMessage += ` (Retry after ${data.retry_after} seconds)`;
          }
  
          errorText.textContent = errorMessage;
          errorMsg.classList.remove("hidden");
        }
  
        // Re-enable submit button
        submitBtn.disabled = false;
        submitBtn.textContent = "Send";
      }
    } catch (error) {
      // Network or other error
      if (errorMsg && errorText) {
        errorText.textContent =
          "Network error. Please check your connection and try again.";
        errorMsg.classList.remove("hidden");
      }
  
      submitBtn.disabled = false;
      submitBtn.textContent = "Send";
    }
  }
  
  function hideMessages() {
    const successMsg = document.getElementById("email-success-message");
    const errorMsg = document.getElementById("email-error-message");
  
    if (successMsg) successMsg.classList.add("hidden");
    if (errorMsg) errorMsg.classList.add("hidden");
  }
  
  // Make functions available globally
  window.closeEmailForm = closeEmailForm;
  window.submitEmailForm = submitEmailForm;