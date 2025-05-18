# Genesis Case: Weather Subscription Service

## Quick Start

1. **Clone the repo:**
   ```sh
   git clone https://github.com/kievzenit/genesis-case.git
   cd genesis-case
   ```

2. **Configure environment:**
   - Create a `.env` file with `WAPP_WEATHER_API_KEY`.

3. **Start with Docker Compose:**
   ```sh
   docker compose up -d
   ```

- Test emails are available at [http://localhost:6569](http://localhost:6569) via the built-in Papercut provider.
- Any email/password can be used for testing.

### Real Email Provider Setup

Set the following in your `.env`:
- `WAPP_EMAIL_SSL=true`
- `WAPP_EMAIL_HOST` (e.g., `smtp.google.com`)
- `WAPP_EMAIL_PORT=465`
- `WAPP_EMAIL_USERNAME`, `WAPP_EMAIL_PASSWORD`, `WAPP_EMAIL_FROM` (for Gmail, `WAPP_EMAIL_FROM` should match `WAPP_EMAIL_USERNAME`).

### Notes

- The `/subscribe` endpoint supports both `application/json` and `application/x-www-form-urlencoded` as per the API specification.