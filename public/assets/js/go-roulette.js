function Rooms (gameName){
    function getSocket(str) {
        return new WebSocket(str);
    }

    this.room_token = "";

    this.msg_socket = null;

    this.ping_socket = null;

    this.player_token = "";

    // next 4 guys is all, that you need to use

    // copy of WebSoket().onclose
    this.onclose = function(event){

    };

    /**
     * redefine this in your code
     * @param data contains 1 if you're first member of room and 2 otherwise
     */
    this.onconnect = function(data){

    };

    /**
     * redefine this in your code
     * @param data contains json message from your game
     */
    this.onmessage = function (data) {

    };

    /**
     * sends json to game
     * @param object json object
     */
    this.send = function (object) {
        $.ajax({
            url: "http://" + location.host + "/api/" + this.room_token + "/action/" + this.player_token + "/",
            data: JSON.stringify(object),
            cache: false,
            type: "POST",
            contentType: "application/json",
        });
    };

    this.connect = function () {
        parent = this;
        $.ajax({
            url: "http://" + location.host + "/new/" + gameName + "/",
            cache: false,
            sync: true,
            success: function (data) {
                parent.room_token = data;
                msg_socket = getSocket("ws://" + location.host + "/ws/" + parent.room_token + "/");
                msg_socket.onclose = onclose;
                msg_socket.onmessage = function (event) {
                    player_token = event.data;
                    ping_socket = getSocket("ws://" + location.host + "/ws/" + parent.room_token + "/" + player_token + "/");
                    ping_socket.onmessage = function () {
                        ping_socket.send(1);
                    };
                    msg_socket.onmessage = function (event){
                        if(event.data == "1" || event.data == "2"){
                            parent.onconnect(event.data);
                        }else{
                            parent.onmessage(event);
                        }
                    };
                };
            }
        });
    };
    return this
}
