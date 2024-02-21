import React, { useState } from "react";
import TextField from "@mui/material/TextField";
import InputAdornment from "@mui/material/InputAdornment";
import IconButton from "@mui/material/IconButton";
import Visibility from "@mui/icons-material/Visibility";
import VisibilityOff from "@mui/icons-material/VisibilityOff";
import { useAuth } from "./AuthContext";

export default function LoginForm() {
  const [showPassword, setShowPassword] = useState(false);
  const [apikey, setApikey] = useState("");
  const { setApiKey, LoginError, setLoginError } = useAuth();

  const handleLogin = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setApiKey(apikey); // Update the global API key
  };

  const handleClickShowPassword = () => {
    setShowPassword(!showPassword);
  };

  const handleMouseDownPassword = (
    event: React.MouseEvent<HTMLButtonElement>
  ) => {
    event.preventDefault();
  };
  const handleApikeyChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setApikey(event.target.value);
    //clear the error message when the user starts typing
    if (LoginError) {
      setLoginError("");
    }
  };

  return (
    <form
      onSubmit={handleLogin}
      style={{
        border: "1px solid #ccc",
        padding: "20px",
        borderRadius: "10px",
        maxWidth: "400px",
        margin: "0 auto",
      }}
    >
      <TextField
        label="API Key"
        variant="outlined"
        onChange={handleApikeyChange}
        error={LoginError ? true : false}
        name="apikey"
        margin="normal"
        type={showPassword ? "text" : "password"}
        required={true}
        fullWidth
        InputProps={{
          endAdornment: (
            <InputAdornment position="end">
              <IconButton
                aria-label="toggle password visibility"
                onClick={handleClickShowPassword}
                onMouseDown={handleMouseDownPassword}
                edge="end"
              >
                {showPassword ? <VisibilityOff /> : <Visibility />}
              </IconButton>
            </InputAdornment>
          ),
        }}
      />
      {LoginError && <p style={{ color: "red" }}>{LoginError}</p>}
      <button className="login-button" type="submit">
        Log In
      </button>

      <hr style={{ margin: "16px 0" }} />
      <button className="signup-button" type="button">
        Create new account
      </button>
    </form>
  );
}
