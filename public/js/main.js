$(function() {
    var FADE_TIME = 150; // ms
    var TYPING_TIMER_LENGTH = 400; // ms
    var COLORS = [
        '#e21400', '#91580f', '#f8a700', '#f78b00',
        '#58dc00', '#287b00', '#a8f07a', '#4ae8c4',
        '#3b88eb', '#3824aa', '#a700ff', '#d300e7'
    ];

    // Initialize variables
    var $window = $(window);
    var $usernameInput = $('.usernameInput'); // Input for username
    var $messages = $('.messages'); // Messages area
    var $inputMessage = $('.inputMessage'); // Input message input box

    var $loginPage = $('.login.page'); // The login page
    var $chatPage = $('.chat.page'); // The chatroom page
    var $loginBtn = $('#submit'); // The chatroom page

    // Prompt for setting a username
    var username;
    var connected = false;
    var typing = false;
    var lastTypingTime;
    var $currentInput = $usernameInput.focus();

    var socket;

    // 添加一条上线通知
    function addParticipantsMessage (data) {
        var message = '';
        if (data.numUsers === 1) {
            message += "there's 1 participant";
        } else {
            message += "there are " + data.numUsers + " participants";
        }
        log(message);
    }

    // 登陆
    function login() {
        var data = {
            username: $("[name='username']").val(),
            passwd: $("[name='passwd']").val()
        };
        $.post("/api/login", data, function (ret) {
            if(ret.code == 200) {
                // 登陆成功
                connedSocket();
                username = $("[name='username']").val();

                $loginPage.fadeOut();
                $chatPage.show();
                $loginPage.off('click');
                $currentInput = $inputMessage.focus();

            } else {
                alert(ret.msg)
            }
            console.log(ret)
        },'json')
    }

    // 发送一条聊天消息
    function sendMessage () {
        var message = $inputMessage.val();
        // 拿到一条消息，并清空输入框
        message = cleanInput(message);
        // if there is a non-empty message and a socket connection
        if (message && connected) {
            $inputMessage.val('');
            addChatMessage({
                username: username,
                message: message
            });
            // tell server to execute 'new message' and send along one parameter
            socket.send(JSON.stringify({Type: "text", message: message, username: username}));
        }
    }

    // 显示一条log消息
    function log (message, options) {
        var $el = $('<li>').addClass('log').text(message);
        addMessageElement($el, options);
    }

    // 添加一条消息，到消息列表
    function addChatMessage (data, options) {
        // Don't fade the message in if there is an 'X was typing'
        var $typingMessages = getTypingMessages(data);
        options = options || {};
        if ($typingMessages.length !== 0) {
            options.fade = false;
            $typingMessages.remove();
        }

        var $usernameDiv = $('<span class="username"/>')
            .text(data.username)
            .css('color', getUsernameColor(data.username));
        var $messageBodyDiv = $('<span class="messageBody">')
            .text(data.message);

        var typingClass = data.typing ? 'typing' : '';
        var $messageDiv = $('<li class="message"/>')
            .data('username', data.username)
            .addClass(typingClass)
            .append($usernameDiv, $messageBodyDiv);

        addMessageElement($messageDiv, options);
    }

    // 添加正在输入动作
    function addChatTyping (data) {
        data.typing = true;
        data.message = '正在输入……';
        addChatMessage(data);
    }

    // 移除正在输入动作
    function removeChatTyping (data) {
        getTypingMessages(data).fadeOut(function () {
            $(this).remove();
        });
    }

    // Adds a message element to the messages and scrolls to the bottom
    // el - The element to add as a message
    // options.fade - If the element should fade-in (default = true)
    // options.prepend - If the element should prepend
    //   all other messages (default = false)
    function addMessageElement (el, options) {
        var $el = $(el);

        // Setup default options
        if (!options) {
            options = {};
        }
        if (typeof options.fade === 'undefined') {
            options.fade = true;
        }
        if (typeof options.prepend === 'undefined') {
            options.prepend = false;
        }

        // Apply options
        if (options.fade) {
            $el.hide().fadeIn(FADE_TIME);
        }
        if (options.prepend) {
            $messages.prepend($el);
        } else {
            $messages.append($el);
        }
        $messages[0].scrollTop = $messages[0].scrollHeight;
    }

    // 清理输入框内容
    function cleanInput (input) {
        return $('<div/>').text(input).html();
    }

    // 更新用户正在输入状态
    function updateTyping () {
        if (connected) {
            if (!typing) {
                typing = true;
                socket.send(JSON.stringify({Type: "typing", username: username}));
            }
            lastTypingTime = (new Date()).getTime();

            setTimeout(function () {
                var typingTimer = (new Date()).getTime();
                var timeDiff = typingTimer - lastTypingTime;
                if (timeDiff >= TYPING_TIMER_LENGTH && typing) {
                    socket.send(JSON.stringify({Type: "stop_typing", username: username}));
                    typing = false;
                }
            }, TYPING_TIMER_LENGTH);
        }
    }

    // Gets the 'X is typing' messages of a user
    function getTypingMessages (data) {
        return $('.typing.message').filter(function (i) {
            return $(this).data('username') === data.username;
        });
    }

    // 获得用户的颜色
    function getUsernameColor (username) {
        // Compute hash code
        var hash = 7;
        for (var i = 0; i < username.length; i++) {
            hash = username.charCodeAt(i) + (hash << 5) - hash;
        }
        // Calculate color
        var index = Math.abs(hash % COLORS.length);
        return COLORS[index];
    }

    // 按键监听
    $window.keydown(function (event) {
        // Auto-focus the current input when a key is typed
        if (!(event.ctrlKey || event.metaKey || event.altKey)) {
            //$currentInput.focus();
        }
        // 回车键登陆
        if (event.which === 13) {
            if (username && connected) {
                sendMessage();
                socket.send(JSON.stringify({Type: "stop_typing", username: username}));
                typing = false;
            } else {
                login();
            }
        }
    });

    // 当消息输入框输入事件时，触发正在输入动作
    $inputMessage.on('input', function() {
        updateTyping();
    });
    $loginBtn.on('click', function () {
        login()
    });

    // Click events

    // 点击输入框时input获得焦点
    $inputMessage.click(function () {
        $inputMessage.focus();
    });

    // Socket events


    // 新用户上线处理
    function newUser(data) {
        connected = true;
        // Display the welcome message
        var message = data.username + " 欢迎来到聊天室.";
        log(message, {prepend: true});
        addParticipantsMessage(data);
    }

    // 登陆成功后连接socket
    function connedSocket() {
        if (window["WebSocket"]) {
            socket = new WebSocket(ws_url);
            socket.onopen = function() {
                connected = true
                console.log("connn websocket")
            };

            socket.onmessage = function(evt) {
                var data = JSON.parse(evt.data);
                if(data.code == 401) {
                    location.href = "/login"
                }
                switch(data.Content.Type) {
                    case "OnlineNotice":
                        newUser({numUsers: 1, })
                        break;
                    case "OfflineNotice":
                        // 用户下线
                        break;
                    default:
                    case "platform":
                        var msg = JSON.parse(data.Content.Data);
                        switch (msg.Type) {
                            default:
                            case "text":
                                addChatMessage({username: msg.username, message: msg.message});
                                break;
                            case 'stop_typing':
                                removeChatTyping({username: msg.username});
                                break;
                            case 'typing':
                                addChatTyping({username: msg.username});
                                break;
                        }
                }
            };

            socket.onerror = function() {
                log("连接异常")
            };

            socket.onclose = function() {
                connected = false;
                $currentInput = $inputMessage.fadeOut();
                $loginPage.focus();
                $loginPage.show();
                $chatPage.hide();
                $loginPage.off('click');

                log("连接已断开")
            };
        } else {
            alert("Your browser does not support WebSockets.");
        }
    }

});