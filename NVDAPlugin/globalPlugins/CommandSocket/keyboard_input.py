from typing import List
import re

# START keys from https://github.com/SeleniumHQ/selenium/blob/63e8543a39778a5fa65ff2149597066abdbb9abb/py/selenium/webdriver/common/keys.py

# Licensed to the Software Freedom Conservancy (SFC) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The SFC licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

"""The Keys implementation."""


class Keys:
	"""Set of special keys codes."""

	NULL = "\ue000"
	CANCEL = "\ue001"  # ^break
	HELP = "\ue002"
	BACKSPACE = "\ue003"
	BACK_SPACE = BACKSPACE
	TAB = "\ue004"
	CLEAR = "\ue005"
	RETURN = "\ue006"
	ENTER = "\ue007"
	SHIFT = "\ue008"
	LEFT_SHIFT = SHIFT
	CONTROL = "\ue009"
	LEFT_CONTROL = CONTROL
	ALT = "\ue00a"
	LEFT_ALT = ALT
	PAUSE = "\ue00b"
	ESCAPE = "\ue00c"
	SPACE = "\ue00d"
	PAGE_UP = "\ue00e"
	PAGE_DOWN = "\ue00f"
	END = "\ue010"
	HOME = "\ue011"
	LEFT = "\ue012"
	ARROW_LEFT = LEFT
	UP = "\ue013"
	ARROW_UP = UP
	RIGHT = "\ue014"
	ARROW_RIGHT = RIGHT
	DOWN = "\ue015"
	ARROW_DOWN = DOWN
	INSERT = "\ue016"
	DELETE = "\ue017"
	SEMICOLON = "\ue018"
	EQUALS = "\ue019"

	NUMPAD0 = "\ue01a"  # number pad keys
	NUMPAD1 = "\ue01b"
	NUMPAD2 = "\ue01c"
	NUMPAD3 = "\ue01d"
	NUMPAD4 = "\ue01e"
	NUMPAD5 = "\ue01f"
	NUMPAD6 = "\ue020"
	NUMPAD7 = "\ue021"
	NUMPAD8 = "\ue022"
	NUMPAD9 = "\ue023"
	MULTIPLY = "\ue024"
	ADD = "\ue025"
	SEPARATOR = "\ue026"
	SUBTRACT = "\ue027"
	DECIMAL = "\ue028"
	DIVIDE = "\ue029"

	F1 = "\ue031"  # function  keys
	F2 = "\ue032"
	F3 = "\ue033"
	F4 = "\ue034"
	F5 = "\ue035"
	F6 = "\ue036"
	F7 = "\ue037"
	F8 = "\ue038"
	F9 = "\ue039"
	F10 = "\ue03a"
	F11 = "\ue03b"
	F12 = "\ue03c"

	META = "\ue03d"
	COMMAND = "\ue03d"
	ZENKAKU_HANKAKU = "\ue040"


# END keys

# Mapping from the unicode key format of webdriver to the NVDA vkCodes
# https://github.com/nvaccess/nvda/blob/master/source/vkCodes.py

modifiers = {
	'NVDA': 'NVDA',
	Keys.SHIFT: 'shift',
	Keys.CONTROL: 'control',
	Keys.ALT: 'alt',
	Keys.INSERT: 'insert',
	Keys.META: 'alt'
}

non_modifiers = {
	Keys.NULL: 'null',
	Keys.CANCEL: "break",
	Keys.HELP: 'help',
	Keys.BACKSPACE: 'backspace',
	Keys.TAB: 'tab',
	Keys.CLEAR: 'clear',
	Keys.RETURN: 'enter',
	Keys.ENTER: 'enter',
	Keys.SHIFT: 'shift',
	Keys.CONTROL: 'control',
	Keys.ALT: 'alt',
	Keys.PAUSE: 'pause',
	Keys.ESCAPE: 'escape',
	Keys.SPACE: 'space',
	Keys.PAGE_UP: 'pageUp',
	Keys.PAGE_DOWN: 'pageDown',
	Keys.END: 'end',
	Keys.HOME: 'home',
	Keys.LEFT: 'leftArrow',
	Keys.UP: 'upArrow',
	Keys.RIGHT: 'rightArrow',
	Keys.DOWN: 'downArrow',
	Keys.INSERT: 'insert',
	Keys.DELETE: 'delete',
	Keys.SEMICOLON: ';',
	Keys.EQUALS: '=',
	Keys.NUMPAD0: 'numpad0',
	Keys.NUMPAD1: 'numpad1',
	Keys.NUMPAD2: 'numpad2',
	Keys.NUMPAD3: 'numpad3',
	Keys.NUMPAD4: 'numpad4',
	Keys.NUMPAD5: 'numpad5',
	Keys.NUMPAD6: 'numpad6',
	Keys.NUMPAD7: 'numpad7',
	Keys.NUMPAD8: 'numpad8',
	Keys.NUMPAD9: 'numpad9',
	Keys.MULTIPLY: 'numpadMultiply',
	Keys.ADD: 'numpadPlus',
	Keys.SEPARATOR: 'numpadSeparator',
	Keys.SUBTRACT: 'numpadMinus',
	Keys.DECIMAL: 'numpadDecimal',
	Keys.DIVIDE: 'numpadDivide',
	Keys.F1: 'f1',
	Keys.F2: 'f2',
	Keys.F3: 'f3',
	Keys.F4: 'f4',
	Keys.F5: 'f5',
	Keys.F6: 'f6',
	Keys.F7: 'f7',
	Keys.F8: 'f8',
	Keys.F9: 'f9',
	Keys.F10: 'f10',
	Keys.F11: 'f11',
	Keys.F12: 'f12',
	Keys.META: 'alt',
	Keys.COMMAND: 'command',
	Keys.ZENKAKU_HANKAKU: 'zenkaku_hankaku',
}


def is_invalid(key: str) -> bool:
	return re.compile(r'^(?:\\\\uE0[A-Z\d]{2}|.)$', re.DOTALL).match(key) is None


def create_gesture_name(keys: List[str]) -> str:
	parts = list()

	for key in keys:
		if is_invalid(key):
			print(f'invalid key \'{key}\'')
			continue

		if key in modifiers.keys():
			parts.append(modifiers[key])
			continue

		if key in non_modifiers.keys():
			parts.append(non_modifiers[key])
			continue

		if str.isalnum(key):
			parts.append(key)
			continue

		print(f'unknown key \'{key}\'')

	if not parts:
		k = ','.join(keys)
		raise Exception(f'could not create gesture using "{k}"')

	return '+'.join(parts)
