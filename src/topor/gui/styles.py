"""Стили и константы GUI — тёмно-синяя тема."""

# Основные цвета
BG_DARK = "#0a1929"        # Тёмно-синий фон
BG_MEDIUM = "#0d2137"      # Чуть светлее
BG_LIGHT = "#122a42"       # Панели / группы
BG_INPUT = "#0f2640"       # Поля ввода
BORDER = "#1a3a5c"         # Границы
BORDER_LIGHT = "#1e4976"   # Светлые границы
TEXT = "#80e5d0"            # Светло-бирюзовый текст
TEXT_DIM = "#5cb8a5"        # Приглушённый бирюзовый
TEXT_BRIGHT = "#a0fff0"     # Яркий бирюзовый
ACCENT = "#00bcd4"         # Акцентный бирюзовый
ACCENT_HOVER = "#00acc1"
ACCENT_PRESS = "#0097a7"
DANGER = "#ef5350"
DANGER_HOVER = "#e53935"
SECONDARY = "#37596e"
SECONDARY_HOVER = "#456d84"
DISABLED_BG = "#1a3050"
DISABLED_TEXT = "#3a6070"
SELECTION = "#00838f"
SUCCESS_ROW = "#0d3b2a"
ERROR_ROW = "#3b1a1a"
SKIP_ROW = "#3b3a0d"

APP_STYLESHEET = f"""
* {{
    color: {TEXT};
    font-size: 13px;
}}

QMainWindow {{
    background-color: {BG_DARK};
}}

QWidget {{
    background-color: {BG_DARK};
}}

QGroupBox {{
    font-weight: bold;
    color: {TEXT_BRIGHT};
    border: 1px solid {BORDER};
    border-radius: 6px;
    margin-top: 14px;
    padding-top: 18px;
    background-color: {BG_MEDIUM};
}}

QGroupBox::title {{
    subcontrol-origin: margin;
    left: 12px;
    padding: 0 6px;
    color: {TEXT_BRIGHT};
}}

QTableWidget {{
    background-color: {BG_MEDIUM};
    alternate-background-color: {BG_LIGHT};
    gridline-color: {BORDER};
    selection-background-color: {SELECTION};
    selection-color: {TEXT_BRIGHT};
    border: 1px solid {BORDER};
    border-radius: 4px;
    color: {TEXT};
}}

QTableWidget::item {{
    padding: 4px;
    color: {TEXT};
}}

QHeaderView::section {{
    background-color: {BG_LIGHT};
    color: {TEXT_BRIGHT};
    padding: 5px 8px;
    border: none;
    border-right: 1px solid {BORDER};
    border-bottom: 1px solid {BORDER};
    font-weight: bold;
}}

QPushButton {{
    background-color: {ACCENT};
    color: {BG_DARK};
    border: none;
    border-radius: 5px;
    padding: 7px 18px;
    min-height: 24px;
    font-size: 13px;
    font-weight: bold;
}}

QPushButton:hover {{
    background-color: {ACCENT_HOVER};
}}

QPushButton:pressed {{
    background-color: {ACCENT_PRESS};
}}

QPushButton:disabled {{
    background-color: {DISABLED_BG};
    color: {DISABLED_TEXT};
}}

QPushButton#dangerButton {{
    background-color: {DANGER};
    color: {BG_DARK};
}}

QPushButton#dangerButton:hover {{
    background-color: {DANGER_HOVER};
}}

QPushButton#secondaryButton {{
    background-color: {SECONDARY};
    color: {TEXT};
}}

QPushButton#secondaryButton:hover {{
    background-color: {SECONDARY_HOVER};
}}

QPushButton#donateButton {{
    background-color: transparent;
    color: {TEXT_BRIGHT};
    font-size: 22px;
    font-weight: bold;
    border: 2px solid {ACCENT};
    border-radius: 18px;
    padding: 0px;
    min-width: 36px;
    max-width: 36px;
    min-height: 36px;
    max-height: 36px;
}}

QPushButton#donateButton:hover {{
    background-color: {ACCENT};
    color: {BG_DARK};
    border-color: {ACCENT};
}}

QTextEdit#logWidget {{
    background-color: {BG_MEDIUM};
    color: {TEXT};
    font-family: "Consolas", "Courier New", monospace;
    font-size: 12px;
    border: 1px solid {BORDER};
    border-radius: 4px;
}}

QProgressBar {{
    border: 1px solid {BORDER};
    border-radius: 5px;
    text-align: center;
    background-color: {BG_LIGHT};
    color: {TEXT_BRIGHT};
    height: 22px;
}}

QProgressBar::chunk {{
    background-color: {ACCENT};
    border-radius: 4px;
}}

QDateTimeEdit, QSpinBox, QLineEdit, QComboBox {{
    padding: 5px 8px;
    border: 1px solid {BORDER};
    border-radius: 4px;
    background-color: {BG_INPUT};
    color: {TEXT};
    min-height: 22px;
    selection-background-color: {SELECTION};
    selection-color: {TEXT_BRIGHT};
}}

QDateTimeEdit:focus, QSpinBox:focus, QLineEdit:focus, QComboBox:focus {{
    border-color: {ACCENT};
}}

QDateTimeEdit::up-button, QSpinBox::up-button,
QDateTimeEdit::down-button, QSpinBox::down-button {{
    background-color: {BG_LIGHT};
    border: 1px solid {BORDER};
    width: 18px;
}}

QDateTimeEdit::up-arrow, QSpinBox::up-arrow,
QDateTimeEdit::down-arrow, QSpinBox::down-arrow {{
    width: 8px;
    height: 8px;
}}

QCheckBox, QRadioButton {{
    spacing: 8px;
    color: {TEXT};
}}

QCheckBox::indicator, QRadioButton::indicator {{
    width: 16px;
    height: 16px;
    border: 2px solid {BORDER_LIGHT};
    background-color: {BG_INPUT};
}}

QCheckBox::indicator:checked, QRadioButton::indicator:checked {{
    background-color: {ACCENT};
    border-color: {ACCENT};
}}

QRadioButton::indicator {{
    border-radius: 10px;
}}

QTabWidget::pane {{
    border: 1px solid {BORDER};
    border-radius: 4px;
    background-color: {BG_MEDIUM};
    top: -1px;
}}

QTabBar::tab {{
    background-color: {BG_LIGHT};
    color: {TEXT_DIM};
    padding: 7px 18px;
    margin-right: 2px;
    border: 1px solid {BORDER};
    border-bottom: none;
    border-top-left-radius: 6px;
    border-top-right-radius: 6px;
}}

QTabBar::tab:selected {{
    background-color: {BG_MEDIUM};
    color: {TEXT_BRIGHT};
    border-bottom: 2px solid {ACCENT};
}}

QTabBar::tab:hover:!selected {{
    background-color: {BG_MEDIUM};
    color: {TEXT};
}}

QSplitter::handle {{
    background-color: {BORDER};
    height: 3px;
}}

QStatusBar {{
    background-color: {BG_MEDIUM};
    color: {TEXT_DIM};
    border-top: 1px solid {BORDER};
}}

QLabel {{
    color: {TEXT};
}}

QMenu {{
    background-color: {BG_MEDIUM};
    color: {TEXT};
    border: 1px solid {BORDER};
}}

QMenu::item:selected {{
    background-color: {SELECTION};
    color: {TEXT_BRIGHT};
}}

QToolTip {{
    background-color: {BG_LIGHT};
    color: {TEXT_BRIGHT};
    border: 1px solid {ACCENT};
    padding: 4px;
    border-radius: 3px;
}}

QScrollBar:vertical {{
    background-color: {BG_DARK};
    width: 12px;
    border: none;
}}

QScrollBar::handle:vertical {{
    background-color: {BORDER_LIGHT};
    border-radius: 5px;
    min-height: 30px;
}}

QScrollBar::handle:vertical:hover {{
    background-color: {ACCENT};
}}

QScrollBar::add-line:vertical, QScrollBar::sub-line:vertical {{
    height: 0px;
}}

QScrollBar:horizontal {{
    background-color: {BG_DARK};
    height: 12px;
    border: none;
}}

QScrollBar::handle:horizontal {{
    background-color: {BORDER_LIGHT};
    border-radius: 5px;
    min-width: 30px;
}}

QScrollBar::handle:horizontal:hover {{
    background-color: {ACCENT};
}}

QScrollBar::add-line:horizontal, QScrollBar::sub-line:horizontal {{
    width: 0px;
}}

QFileDialog {{
    background-color: {BG_DARK};
}}

QMessageBox {{
    background-color: {BG_DARK};
}}

QCalendarWidget {{
    background-color: {BG_MEDIUM};
    color: {TEXT};
}}
"""
