module.exports = {
  apps: [
    {
      name: "sphinx-tribes-backend",
      script: "./sphinx-tribes",
      cwd: "/workspaces/sphinx-tribes",
      instances: 1,
      autorestart: true,
      watch: false,
      max_memory_restart: "1G",
      env: {
        NODE_ENV: "development",
        RESTART: "true",
        REBUILD_COMMAND: "go build -o sphinx-tribes",
      },
    },
    {
      name: "sphinx-tribes-frontend",
      script: "yarn",
      args: "start:codespace",
      cwd: "/workspaces/sphinx-tribes-frontend",
      instances: 1,
      autorestart: true,
      watch: false,
      max_memory_restart: "1G",
      env: {
        INSTALL_COMMAND: "yarn install",
      },
    },
  ],
};
