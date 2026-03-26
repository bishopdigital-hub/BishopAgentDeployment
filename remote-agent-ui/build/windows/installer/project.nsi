Unicode true

!include "wails_tools.nsh"

VIProductVersion "${INFO_PRODUCTVERSION}.0"
VIFileVersion    "${INFO_PRODUCTVERSION}.0"

VIAddVersionKey "CompanyName"     "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion"  "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion"     "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright"  "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName"     "${INFO_PRODUCTNAME}"

ManifestDPIAware true

!include "MUI.nsh"
!include "nsDialogs.nsh"
!include "LogicLib.nsh"

!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"
!define MUI_FINISHPAGE_NOAUTOCLOSE
!define MUI_ABORTWARNING

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY

Page custom AgentNamePage AgentNamePageLeave

!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "English"

Name "${INFO_PRODUCTNAME}"
OutFile "..\..\bin\${INFO_PROJECTNAME}-${ARCH}-installer.exe"
InstallDir "$PROGRAMFILES64\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}"
ShowInstDetails show

Var AgentName
Var AgentNameEdit
Var CustomerName
Var CustomerNameEdit
Var UniqueID
Var UniqueIDEdit
Var DiscordChannelID
Var DiscordChannelIDEdit

Function .onInit
FunctionEnd

Function AgentNamePage
    !insertmacro MUI_HEADER_TEXT "Agent Configuration" "Enter deployment details for this agent."
    nsDialogs::Create 1018
    Pop $0
    ${If} $0 == error
        Abort
    ${EndIf}

    ${NSD_CreateLabel} 0 0 100% 12u "Friendly Name (e.g., Server-XP-01):"
    Pop $0
    ${NSD_CreateText} 0 13u 100% 12u "Remote-Agent"
    Pop $AgentNameEdit

    ${NSD_CreateLabel} 0 30u 100% 12u "Customer/Company Name:"
    Pop $0
    ${NSD_CreateText} 0 43u 100% 12u "Private-Client"
    Pop $CustomerNameEdit

    ${NSD_CreateLabel} 0 60u 100% 12u "Unique ID / Discord Handle:"
    Pop $0
    ${NSD_CreateText} 0 73u 100% 12u "nexus-agent-001"
    Pop $UniqueIDEdit

    ${NSD_CreateLabel} 0 90u 100% 12u "Discord Channel ID (Tactical Link):"
    Pop $0
    ${NSD_CreateText} 0 103u 100% 12u "1485020950457094398"
    Pop $DiscordChannelIDEdit

    nsDialogs::Show
FunctionEnd

Function AgentNamePageLeave
    ${NSD_GetText} $AgentNameEdit $AgentName
    ${NSD_GetText} $CustomerNameEdit $CustomerName
    ${NSD_GetText} $UniqueIDEdit $UniqueID
    ${NSD_GetText} $DiscordChannelIDEdit $DiscordChannelID
FunctionEnd

# Note: Bypassing wails.checkArchitecture for Server 2012 support
# !insertmacro wails.checkArchitecture

Section
    !insertmacro wails.setShellContext
    !insertmacro wails.webview2runtime

    SetOutPath $INSTDIR
    !insertmacro wails.files
    
    # Simple JSON injection for config.json
    # We overwrite the default config.json with one containing the user's chosen details
    FileOpen $0 "$INSTDIR\config.json" w
    FileWrite $0 '{"agent_id": "$UniqueID", "agent_name": "$AgentName", "company": "$CustomerName", "role": "server", "nexus_url": "http://41.193.213.128:9500", "discord_channel_id": "$DiscordChannelID"}'
    FileClose $0

    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
    CreateShortCut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"

    !insertmacro wails.associateFiles
    !insertmacro wails.associateCustomProtocols
    !insertmacro wails.writeUninstaller
SectionEnd

Section "uninstall"
    !insertmacro wails.setShellContext
    RMDir /r "$AppData\${PRODUCT_EXECUTABLE}"
    RMDir /r $INSTDIR
    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"
    !insertmacro wails.deleteUninstaller
SectionEnd
