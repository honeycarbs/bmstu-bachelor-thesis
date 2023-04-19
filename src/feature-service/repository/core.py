from abc import ABC, abstractmethod


class Repository(ABC):
    @abstractmethod
    def create(self, entity):
        pass
