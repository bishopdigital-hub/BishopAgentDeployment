import './app.css';
import logo from './assets/images/logo-bishop.svg';
import { Greet, RunScan, RunAudit, RunExploit } from '../wailsjs/go/main/App';
// @ts-ignore
import { Login } from '../wailsjs/go/main/App';

// UI State
let currentRole = 'server';

// Initialize UI
document.getElementById('logo').src = logo;
document.getElementById('login-logo').src = logo;
const consoleOutput = document.getElementById('console-output');

// Auth Logic
window.verifyLogin = () => {
    const password = document.getElementById('password-input').value;
    const errorMsg = document.getElementById('login-error');
    
    logToConsole('Authenticating Operator...', 'info');
    
    Login(password).then(success => {
        if (success) {
            logToConsole('Tactical Link Established. Access Granted.', 'success');
            document.getElementById('login-overlay').style.display = 'none';
            document.getElementById('app').style.display = 'block';
            sessionStorage.setItem('bishop_auth', 'true');
        } else {
            logToConsole('Authentication Failed: Invalid Operator Credentials.', 'error');
            errorMsg.style.display = 'block';
            document.getElementById('password-input').value = '';
        }
    }).catch(err => {
        logToConsole(`Authentication System Error: ${err}`, 'error');
    });
};

// Check existing session
if (sessionStorage.getItem('bishop_auth') === 'true') {
    document.getElementById('login-overlay').style.display = 'none';
    document.getElementById('app').style.display = 'block';
}

// Helper to log to virtual console
function logToConsole(message, type = 'info') {
    const timestamp = new Date().toLocaleTimeString();
    const style = type === 'error' ? 'color: #ef4444' : (type === 'success' ? 'color: #10b981' : '');
    consoleOutput.innerHTML += `<span style="opacity: 0.5">[${timestamp}]</span> <span style="${style}">${message}</span>\n`;
    consoleOutput.scrollTop = consoleOutput.scrollHeight;
}

// Role Selection logic
window.setRole = (role) => {
    currentRole = role;
    document.getElementById('role-server').classList.toggle('active', role === 'server');
    document.getElementById('role-workstation').classList.toggle('active', role === 'workstation');
    
    logToConsole(`Sector Role initialized: ${role.toUpperCase()}`, 'info');
};

// Scan Logic
window.runScan = () => {
    logToConsole('Initiating Network Topology Scan...', 'info');
    RunScan('127.0.0.1').then(res => {
        logToConsole('Scan Complete. Topology mapping synchronized.', 'success');
        logToConsole(res);
    }).catch(err => {
        logToConsole(`Scan Failed: ${err}`, 'error');
    });
};

// Audit Logic
window.runAudit = () => {
    logToConsole('Initiating System Security Audit (Lynis)...', 'info');
    RunAudit().then(res => {
        logToConsole('Audit Complete. Security telemetry captured.', 'success');
        logToConsole(res);
    }).catch(err => {
        logToConsole(`Audit Failed: ${err}`, 'error');
    });
};

// Exploit Logic
window.runExploit = () => {
    logToConsole('Initiating Autonomous Exploitation Sequence (MSF)...', 'info');
    RunExploit('127.0.0.1').then(res => {
        logToConsole(res, 'success');
    }).catch(err => {
        logToConsole(`Exploit Failed: ${err}`, 'error');
    });
};

// Service Logic
window.launchService = (url) => {
    logToConsole(`Launching service gateway: ${url}`, 'info');
    window.go.main.App.LaunchService(url);
};

// Console controls
window.clearOutput = () => {
    consoleOutput.innerHTML = '';
};

// Periodic Status Polling
setInterval(() => {
    ['n8n', 'webui'].forEach(service => {
        window.go.main.App.GetServiceStatus(service).then(status => {
            const marker = document.getElementById(`status-${service}`);
            if (marker) {
                marker.innerText = status.toUpperCase();
                marker.className = `status-marker ${status}`;
            }
        });
    });
}, 10000);

// Simulate Nexus connectivity for now
setTimeout(() => {
    const badge = document.getElementById('status-badge');
    badge.innerText = 'Nexus Synchronized';
    badge.classList.remove('offline');
    badge.classList.add('online');
    logToConsole('Nexus Synchronization established.', 'success');
}, 2000);
