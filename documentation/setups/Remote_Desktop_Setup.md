# xRDP Setup on Ubuntu with GNOME Desktop

This guide provides step-by-step instructions to set up xRDP on Ubuntu using `gdm3`, `gnome-session`, and the `ubuntu-desktop` environment. This setup allows you to connect to your Ubuntu machine remotely using the full GNOME desktop environment.

## Prerequisites

- Ubuntu 20.04 or later.
- Sudo privileges.

## Steps

### 1. Update the System

Before starting the installation, update your system to ensure all packages are up to date.

```shell
sudo apt update
sudo apt upgrade -y
sudo reboot
```

### 2. Install xRDP

Install xRDP and the required dependencies.

```shell
sudo apt install xrdp -y
```

### 3. Install GNOME Desktop Environment

Ensure the GNOME desktop environment (`ubuntu-desktop`) is installed. This will also install `gdm3` if it is not already installed.

```shell
sudo apt install ubuntu-desktop gnome-session -y
```

### 4. Configure xRDP to Use GNOME

#### 4.1 Modify the xRDP Startup Script

Edit the xRDP startup script to start the GNOME desktop environment when connecting via RDP.

```shell
sudo nano /etc/xrdp/startwm.sh
```

Replace the contents of `startwm.sh` with the following:

```sh
#!/bin/sh
# xrdp X session start script

if test -r /etc/profile; then
    . /etc/profile
fi

if test -r ~/.profile; then
    . ~/.profile
fi

unset DBUS_SESSION_BUS_ADDRESS
unset XDG_RUNTIME_DIR

# Start GNOME session
export GNOME_SHELL_SESSION_MODE=ubuntu
export XDG_CURRENT_DESKTOP=ubuntu:GNOME
export XDG_SESSION_DESKTOP=ubuntu
exec /usr/bin/gnome-session
```

Save the file and exit (`Ctrl + O`, `Enter`, `Ctrl + X`).

#### 4.2 Configure gdm3 to Work with xRDP

By default, `gdm3` can cause issues with xRDP. You need to configure it to allow xRDP sessions:

1. **Edit the Custom Configuration for gdm3**:

   ```shell
   sudo nano /etc/gdm3/custom.conf
   ```

2. **Uncomment and Set the Wayland Enable Line to False**:

   Find the line `[daemon]` and make sure the `WaylandEnable` line is set to `false`:

   ```ini
   [daemon]
   WaylandEnable=false
   ```

3. **Save and Exit** (`Ctrl + O`, `Enter`, `Ctrl + X`).

### 5. Restart Services

Restart the `xrdp` service to apply the changes:

```shell
sudo systemctl restart xrdp
```

Restart the `gdm3` service:

```shell
sudo systemctl restart gdm3
```

### 6. Allow RDP Through the Firewall

If you have `ufw` (Uncomplicated Firewall) enabled, you need to allow RDP traffic:

```shell
sudo ufw allow 3389/tcp
```

### 7. Connect to the Ubuntu Machine via RDP

You can now connect to your Ubuntu machine using an RDP client (like the built-in Remote Desktop Connection on Windows, or Remmina on Linux).

1. **Open your RDP client**.
2. **Enter the IP address** of your Ubuntu machine.
3. **Log in with your Ubuntu credentials**.

You should now be able to access your Ubuntu machine with the full GNOME desktop environment over RDP.

## Troubleshooting

### Black Screen After Login

If you encounter a black screen after logging in:

1. Ensure that `Wayland` is disabled in `/etc/gdm3/custom.conf`.
2. Check that the `startwm.sh` script is correctly configured as shown above.
3. Restart `xrdp` and `gdm3` services.

### Failed to Start Session

If you see a "Failed to start session" error:

1. Ensure that the correct desktop environment is set in the `startwm.sh` script.
2. Make sure all GNOME packages are properly installed.
3. Reboot the system to ensure all configurations are applied.
