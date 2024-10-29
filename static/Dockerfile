FROM nginx

# Copy the custom NGINX configuration if needed
COPY nginx.conf /etc/nginx/

# Copy the HTML, CSS, JavaScript, and icon files to the NGINX directory
COPY index.html /usr/share/nginx/html/
COPY styles.css /usr/share/nginx/html/
COPY script.js /usr/share/nginx/html/
COPY affirmationicon.jpg /usr/share/nginx/html/