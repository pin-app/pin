# Frontend Setup with Authentication and Dev Mode

This document explains how to set up and use the Pin mobile app with the new authentication system and dev mode.

## Features

- **OAuth Authentication**: Google and Apple OAuth integration
- **Dev Mode**: Bypass authentication for development and testing
- **Session Management**: Persistent login sessions
- **API Integration**: Full integration with the backend API
- **Environment Configuration**: Easy switching between dev and production

## Setup

### 1. Install Dependencies

```bash
cd mobile
npm install
```

### 2. Environment Configuration

The app uses environment variables for configuration. You can set these in your shell or create a `.env` file:

```bash
# API Configuration
export EXPO_PUBLIC_API_BASE=http://localhost:8080
export EXPO_PUBLIC_DEV_MODE=true

# OAuth Configuration (optional for dev mode)
export EXPO_PUBLIC_GOOGLE_CLIENT_ID=your_google_client_id
export EXPO_PUBLIC_APPLE_CLIENT_ID=your_apple_client_id
```

### 3. Start the App

#### Development Mode (with dev mode enabled)
```bash
# macOS
npm run start:mac

# Linux
npm run start:linux

# Windows
npm run start:win
```

#### Production Mode
```bash
npm run start:prod
```

## Usage

### Authentication Flow

1. **First Launch**: App shows authentication screen
2. **OAuth Login**: Tap "Continue with Google" or "Continue with Apple"
3. **Dev Mode**: Toggle dev mode and enter a user ID for testing
4. **Session Persistence**: Login state is saved and restored on app restart

### Dev Mode

Dev mode allows you to bypass OAuth authentication for development:

1. **Enable Dev Mode**: Toggle the switch in the auth screen or profile settings
2. **Set Dev User**: Enter a user ID to use for API requests
3. **API Requests**: All API requests will include the `X-Dev-User-ID` header
4. **Visual Indicators**: Dev mode status is shown throughout the app

### API Integration

The app automatically handles:
- **Authentication Headers**: Adds Bearer tokens or dev user headers
- **Error Handling**: Graceful fallback to mock data if API fails
- **Loading States**: Shows loading indicators during API calls
- **Offline Support**: Falls back to cached/mock data when offline

## Development

### File Structure

```
src/
├── contexts/
│   └── AuthContext.tsx          # Authentication state management
├── services/
│   └── api.ts                   # API service with auth integration
├── screens/
│   ├── Auth/
│   │   └── index.tsx            # Authentication screen
│   ├── Feed/
│   │   └── index.tsx            # Feed with real API data
│   ├── Map/
│   │   └── HealthCheck.tsx      # Health check with dev mode info
│   └── Profile/
│       └── index.tsx            # Profile with dev settings
├── components/
│   └── DevModeSettings.tsx      # Dev mode configuration modal
└── config/
    └── environment.ts           # Environment configuration
```

### Key Components

#### AuthContext
- Manages authentication state
- Handles login/logout
- Provides dev mode controls
- Persists state in AsyncStorage

#### API Service
- Centralized API communication
- Automatic authentication header management
- Dev mode support
- Error handling and fallbacks

#### DevModeSettings
- Toggle dev mode on/off
- Set dev user ID
- View current configuration
- Accessible from profile screen

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `EXPO_PUBLIC_API_BASE` | Backend API URL | `http://localhost:8080` |
| `EXPO_PUBLIC_DEV_MODE` | Enable dev mode | `false` |
| `EXPO_PUBLIC_GOOGLE_CLIENT_ID` | Google OAuth client ID | - |
| `EXPO_PUBLIC_APPLE_CLIENT_ID` | Apple OAuth client ID | - |

## Backend Integration

The frontend is designed to work with the Pin backend API:

### Required Backend Features
- OAuth endpoints (`/api/auth/google`, `/api/auth/apple`)
- Session management
- Dev mode support
- Protected endpoints with proper authentication

### API Endpoints Used
- `GET /health` - Health check
- `GET /api/posts` - Fetch posts
- `GET /api/users/{id}` - Get user info
- `POST /api/auth/logout` - Logout

## Troubleshooting

### Common Issues

1. **API Connection Failed**
   - Check `EXPO_PUBLIC_API_BASE` is correct
   - Ensure backend is running
   - Check network connectivity

2. **Dev Mode Not Working**
   - Verify `EXPO_PUBLIC_DEV_MODE=true`
   - Check backend has `DEV_MODE=true`
   - Restart the app after changing settings

3. **OAuth Not Working**
   - Check OAuth client IDs are set
   - Verify redirect URLs are configured
   - Check backend OAuth configuration

### Debug Mode

Enable debug logging by setting:
```bash
export EXPO_PUBLIC_DEBUG=true
```

This will show additional console logs for debugging API calls and authentication flow.

## Production Deployment

For production deployment:

1. Set `EXPO_PUBLIC_DEV_MODE=false`
2. Configure proper OAuth client IDs
3. Set production API base URL
4. Test OAuth flow thoroughly
5. Remove dev mode UI elements if desired

## Security Notes

- Dev mode should never be enabled in production
- OAuth client secrets should be kept secure
- Session tokens are stored securely in AsyncStorage
- API requests use HTTPS in production
