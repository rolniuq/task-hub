package desktop

import (
	"context"
	"fmt"
	"taskhub/internal/app"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DesktopApp struct {
	fyneApp     fyne.App
	authService *app.AuthService
	mainWindow  fyne.Window
	currentUser *app.LoginResponse
}

func NewApp(authService *app.AuthService) *DesktopApp {
	fyneApp := fyneApp.New()
	fyneApp.Settings().SetTheme(&CustomTheme{})

	return &DesktopApp{
		fyneApp:     fyneApp,
		authService: authService,
	}
}

func RunDesktopApp(desktopApp *DesktopApp) {
	desktopApp.createMainWindow()
	desktopApp.showLoginScreen()
	desktopApp.fyneApp.Run()
}

func (d *DesktopApp) createMainWindow() {
	d.mainWindow = d.fyneApp.NewWindow("Task Hub")
	d.mainWindow.Resize(fyne.NewSize(450, 700))
	d.mainWindow.CenterOnScreen()
	d.mainWindow.SetFixedSize(false)
}

func (d *DesktopApp) showLoginScreen() {
	loginForm := d.createLoginForm()
	d.mainWindow.SetContent(loginForm)
	d.mainWindow.Show()
}

func (d *DesktopApp) createLoginForm() *fyne.Container {
	// App title
	titleLabel := widget.NewLabelWithStyle("Task Hub", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	subtitleLabel := widget.NewLabelWithStyle("Welcome back! Please login to continue.", fyne.TextAlignCenter, fyne.TextStyle{})
	subtitleLabel.Wrapping = fyne.TextWrapWord

	// Email entry with icon
	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Enter your email")
	emailEntry.Resize(fyne.NewSize(300, 40))

	// Password entry with icon
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter your password")
	passwordEntry.Resize(fyne.NewSize(300, 40))

	// Styled login button
	loginBtn := widget.NewButtonWithIcon("Login", theme.LoginIcon(), func() {
		d.handleLogin(emailEntry.Text, passwordEntry.Text)
	})
	loginBtn.Importance = widget.HighImportance

	// Register link
	registerLabel := widget.NewHyperlink("Don't have an account? Register", nil)
	registerLabel.OnTapped = func() {
		d.showRegisterScreen()
	}

	// Form layout with better spacing
	form := container.NewVBox(
		container.NewVBox(
			titleLabel,
			subtitleLabel,
		),
		container.NewPadded(
			container.NewVBox(
				widget.NewLabel("Email"),
				emailEntry,
				widget.NewLabel("Password"),
				passwordEntry,
			),
		),
		container.NewPadded(loginBtn),
		container.NewCenter(registerLabel),
	)

	// Wrap in a card for better visual appeal
	card := widget.NewCard("", "", form)
	card.Resize(fyne.NewSize(350, 400))

	return container.NewCenter(card)
}

func (d *DesktopApp) handleLogin(email, password string) {
	if email == "" || password == "" {
		dialog.ShowError(fmt.Errorf("email and password are required"), d.mainWindow)
		return
	}

	ctx := context.Background()
	req := &app.LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := d.authService.Login(ctx, req)
	if err != nil {
		dialog.ShowError(fmt.Errorf("login failed: %v", err), d.mainWindow)
		return
	}

	d.currentUser = resp
	d.showDashboard()
}

func (d *DesktopApp) showRegisterScreen() {
	// App title
	titleLabel := widget.NewLabelWithStyle("Create Account", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	subtitleLabel := widget.NewLabelWithStyle("Join Task Hub today!", fyne.TextAlignCenter, fyne.TextStyle{})
	subtitleLabel.Wrapping = fyne.TextWrapWord

	// Form entries
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter your full name")
	nameEntry.Resize(fyne.NewSize(300, 40))

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Enter your email")
	emailEntry.Resize(fyne.NewSize(300, 40))

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter your password")
	passwordEntry.Resize(fyne.NewSize(300, 40))

	confirmPasswordEntry := widget.NewPasswordEntry()
	confirmPasswordEntry.SetPlaceHolder("Confirm your password")
	confirmPasswordEntry.Resize(fyne.NewSize(300, 40))

	// Styled register button
	registerBtn := widget.NewButtonWithIcon("Create Account", theme.ConfirmIcon(), func() {
		d.handleRegister(nameEntry.Text, emailEntry.Text, passwordEntry.Text, confirmPasswordEntry.Text)
	})
	registerBtn.Importance = widget.HighImportance

	// Back link
	backLabel := widget.NewHyperlink("Already have an account? Login", nil)
	backLabel.OnTapped = func() {
		d.showLoginScreen()
	}

	// Form layout with better spacing
	form := container.NewVBox(
		container.NewVBox(
			titleLabel,
			subtitleLabel,
		),
		container.NewPadded(
			container.NewVBox(
				widget.NewLabel("Full Name"),
				nameEntry,
				widget.NewLabel("Email"),
				emailEntry,
				widget.NewLabel("Password"),
				passwordEntry,
				widget.NewLabel("Confirm Password"),
				confirmPasswordEntry,
			),
		),
		container.NewPadded(registerBtn),
		container.NewCenter(backLabel),
	)

	// Wrap in a card for better visual appeal
	card := widget.NewCard("", "", form)
	card.Resize(fyne.NewSize(350, 500))

	d.mainWindow.SetContent(container.NewCenter(card))
}

func (d *DesktopApp) handleRegister(name, email, password, confirmPassword string) {
	if name == "" || email == "" || password == "" {
		dialog.ShowError(fmt.Errorf("all fields are required"), d.mainWindow)
		return
	}

	if password != confirmPassword {
		dialog.ShowError(fmt.Errorf("passwords do not match"), d.mainWindow)
		return
	}

	ctx := context.Background()
	req := &app.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	}

	_, err := d.authService.Register(ctx, req)
	if err != nil {
		dialog.ShowError(fmt.Errorf("registration failed: %v", err), d.mainWindow)
		return
	}

	dialog.ShowInformation("Success", "Registration successful! Please login.", d.mainWindow)
	d.showLoginScreen()
}

func (d *DesktopApp) showDashboard() {
	// Welcome header
	welcomeLabel := widget.NewLabelWithStyle(fmt.Sprintf("Welcome, %s! 👋", d.currentUser.User.Name), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// User info card
	userInfo := container.NewVBox(
		widget.NewLabelWithStyle("📧 Email", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(d.currentUser.User.Email),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("👤 Account Status", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Active"),
	)

	userCard := widget.NewCard("User Information", "", userInfo)

	// Quick actions section
	actionsLabel := widget.NewLabelWithStyle("Quick Actions", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	logoutBtn := widget.NewButtonWithIcon("Logout", theme.LogoutIcon(), func() {
		d.handleLogout()
	})
	logoutBtn.Importance = widget.MediumImportance

	// Dashboard content with better layout
	content := container.NewVBox(
		container.NewPadded(welcomeLabel),
		container.NewPadded(userCard),
		container.NewVBox(
			actionsLabel,
			logoutBtn,
		),
	)

	// Main content with padding
	mainContent := container.NewPadded(content)
	d.mainWindow.SetContent(mainContent)
}

func (d *DesktopApp) handleLogout() {
	d.currentUser = nil
	d.showLoginScreen()
}
