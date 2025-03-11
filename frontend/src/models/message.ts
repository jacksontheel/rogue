interface BaseMessage {
  type: string;
}

interface UserIdMessage extends BaseMessage {
  type: "userId";
  userId: string;
}

interface PositionMessage extends BaseMessage {
  type: "position";
  data: {
    x: number;
    y: number;
  };
}

interface PublishPositionMessage extends PositionMessage {
  userId: string;
}

interface PublishPlayerExitMessage extends BaseMessage {
  type: "playerExit",
  userId: string
}

export type incomingMessage = UserIdMessage | PublishPositionMessage | PublishPlayerExitMessage;
export type outgoingMessage = PositionMessage;
