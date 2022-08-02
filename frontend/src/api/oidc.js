import jwt_decode from 'jwt-decode';
import { v4 as uuidv4 } from 'uuid';

/** Determines if application is authenticated */
export const getAuthenticated = async () => {
  console.log("getAuthenticated()")
  const session_code = window.localStorage.getItem('session_code');
  const provider = window.localStorage.getItem('provider');
  if (session_code === null) {
    return false;
  } else if (provider === null) {
    return false;
  } else {
  const response = await fetch(`http://localhost:8090/authenticated?provider=${provider}`, {
    headers: {
      Authorization: `Bearer ${session_code}`,
    },
  });
  if (!response.ok) {
    throw new Error();
  }
  const { status } = await response.json();
  if (status === 'authenticated') {
    return true;
  } else {
    try {
      await refresh();
      return true;
    } catch (err) {
      logout();
      throw new Error();
    }
  }
  }
}

/** Prepares application for login and loads login screen  */
export const loadLoginScreen = () => {
  const state = uuidv4();
  const nonce = uuidv4();
  window.localStorage.setItem('state', state);
  window.localStorage.setItem('nonce', nonce);
  const provider = window.localStorage.getItem('provider');
  window.location.assign(`http://localhost:8090/login-screen?state=${state}&nonce=${nonce}&provider=${provider}`);
}

/** Validates returned state against persisted state */
export const validateState = checkState => {
  const state = window.localStorage.getItem('state');
  window.localStorage.removeItem('state');
  return checkState === state;
}


/** Validates returned nonce against persisted nonce */
const validateNonce = checkNonce => {
  const nonce = window.localStorage.getItem('nonce');
  window.localStorage.removeItem('nonce');
  return checkNonce === nonce;
}

/** Exchange code for tokens  */
export const login = async code => {
  const provider = window.localStorage.getItem('provider');
  const response = await fetch(
    `http://localhost:8090/get-session-code?provider=${provider}`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ code }), 
    },
  );
  if (!response.ok) {
    return false;
  }
  const { session_code, nonce } = await response.json();
  const validNonce = validateNonce(nonce);
  if (!validNonce) return false;
  window.localStorage.setItem('session_code', session_code);
  return true;
}

export const refresh = async () => {
  const old_session_code = window.localStorage.getItem('session_code');
  const provider = window.localStorage.getItem('provider');
  const response = await fetch(
    `http://localhost:8090/refresh-tokens?provider=${provider}`,
    {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${old_session_code}`,
      }, 
    },
  );
  if (!response.ok) {
    throw new Error();
  }
  const { session_code } = await response.json();
  window.localStorage.setItem('session_code', session_code);
}

export const logout = () => {
  window.localStorage.removeItem('session_code');
  window.localStorage.removeItem('provider');
  // window.localStorage.removeItem('refresh_token');
  window.location.reload();
}

export const logoutSession = async () => {
  const session_code = window.localStorage.getItem('session_code');
  const provider = window.localStorage.getItem('provider');
  const response = await fetch(
    `http://localhost:8090/logout?provider=${provider}`,
    {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${session_code}`,
      },
    },
  );
  if (!response.ok) {
    throw new Error();
  }
  const logoutResponse = await response.text();
  console.log(logoutResponse);
  window.localStorage.removeItem('session_code');
  window.localStorage.removeItem('provider');
}
